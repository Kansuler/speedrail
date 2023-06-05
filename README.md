# Speedrail
![Speedrail](./speedrail.png)

## Introduction
Speedrail is a lib that will help you compose a plan of strategies that will execute in order. It's super handy for
anyone that builds a modular system with a lot of conditions and unique edge-cases. The lib rely heavily on generics to
let you define your own strategy signatures with custom data models, service containers and conditions.

## Example
A fast example that shows the readability of Speedrail when used in your project.

```go
func main() {
    container := Container {
        DB: sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
    }
	
	var model Model
	
    plan := speedrail.Plan(
		strategy.ParseRequest,
		speedrail.If(
			speedrail.Or(
				speedrail.Not(condition.HasUsername),
				speedrail.Not(condition.HasPassword),
			),
			speedrail.ThrowError(speedrail.NewError(errors.New("request missing username or password"), http.StatusBadRequest, "missing username or password")),
		),
		strategy.InsertUserToDB,
		strategy.CreateToken,
	)
	
	ctx, model, err := plan.Execute(context.Background(), container, model)
	if err != nil {
        // handle error
    }
	
	fmt.Println(model)
}
```

### Define strategy signature
All strategies that you build will have to follow the `speedrail.Strategy[C, M any]` signature. `C` is short for
`Container` and this struct will hold all clients, database connection pools and other services that you want strategies
to have access to. Strategies will have access, but should not modify the container.

`M` is short for `Model` and this is the data model that you want to pass to the strategy. The model can be anything
that you define. It can be a struct, a map or a slice. Refrain from using pointers in the Model. Strategies will receive
a copy of the model and may mutate its data. When a strategy is done, it should return the model so that it gets passed
along to the next strategy.

#### Container example

Here is an example of a container with a `sql.DB` instance and a `http.Client` instance. These will be reachable in each
strategy and condition that you will create.

```go
type Container struct {
    DB *sql.DB // Hold an open DB connection pool
    Client *http.Client // Hold an instantiated http client, with settings.
}
```

#### Model example
In a model you should define fields with mutable data that is relevant for your strategies. The model can be anything
you want it to be.

```go
type Model struct {
    ID string json:"-"
    UserName string `json:"username"`
    Password string `json:"-"`
    Email string `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}
```

#### Strategy signature
A strategy signature may be used for convenience. You can use it to type check that your strategies are following the
common signature during compile time.

```go
type Signature = speedrail.Strategy[Container, Model]

// Strategy example that follows the signature. Note that it retrieves Container and Model, and return Model just as
// defined above.
func ExampleStrategy(ctx context.Context, container Container, model Model) (context.Context, Model, speedrail.Error) {
    // ... do things
    return ctx, model, nil
}

// Type check the signature so that it follows Signature. This will display error in your IDE, or fail to compile should
// the strategy not follow the signature.
var _ Signature = ExampleStrategy
```

### Define a strategy
Now we can start building our strategies! Think of a strategy as a modular function that can be executed given the right
conditions.

#### Strategy that mutates the data model
A simple strategy that mutates the data model. In this case it changes UserName and then returns the mutated model.

```go
func SetUserName(ctx context.Context, container Container, model Model) (context.Context, Model, speedrail.Error) {
    model.UserName = "John Doe"
    return ctx, model, nil
}

// Type checking so that strategy follows our defined signature
var _ Signature = SetUserName
```

#### Strategy that returns an error
A simple strategy that mutates the data model. In this case it changes UserName and then returns the mutated model.

```go
// Example where an error has occurred.
func GetUser(ctx context.Context, container Container, model Model) (context.Context, Model, speedrail.Error) {
    user, err := container.Client.GetUserName(model.UserName)
	if err != nil {
		return ctx, model, speedrail.NewError(err, http.StatusInternalServerError, "something went wrong calling on client")
    }
	model.ID = user.ID
    return ctx, model, nil
}

// Type checking so that strategy follows our defined signature
var _ Signature = GetUser
```

#### Strategy that uses a service from the container
A strategy that uses a service from the container. In this case it uses the `sql.DB` service to make an insert query to
database, and then scan in the result to the model. The mutated model is then returned. If an error occurs, it will be
returned as a speedrail.Error. You may create your own error types, as long as it follows the speedrail.Error interface.

```go
func InsertUserToDatabase(ctx context.Context, container Container, model Model) (context.Context, Model, speedrail.Error) {
    query := container.DB.QueryRow("INSERT INTO users (username, email, created_at) VALUES (?, ?, ?) RETURNING id", model.ExternalData.UserName, model.ExternalData.Email, model.ExternalData.CreatedAt)
    err := query.Scan(&model.ID)
	if err != nil {
		// Create a new error and return it, if you'd wish you can create your own error type that implements
		// speedrail.Error interface.
        return ctx, model, speedrail.NewError(err, http.StatusInternalServerError, "some error")
    }
	
    return ctx, model, nil
}

// Type checking so that strategy follows our defined signature
var _ Signature = InsertUserToDatabase
```

### Plan
A plan is a list of strategies that compose a plan of execution. All strategies that are part of the same plan must have
the same strategy signature. You may define different plans with different strategy signatures, depending on your needs.

```go
func main() {
    // Instantiate a container
    container := Container{
        DB: sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
        Client: http.DefaultClient,
    }
	
    // Instantiate a model
    model := DataModel{
        CreatedAt: time.Now(),
        Email: "",
    }
	
    // Create a plan of strategies
    plan := speedrail.Plan(
        SetUserName,
        InsertUserToDatabase,
    )
	
    var err speedrail.Error
    // Execute the plan with given context, container and model. The strategies will execute in order, and return error
	// if any of them fails.
    ctx, model, err = plan.Execute(context.Background(), container, model)
    if err != nil {
        panic(err)
    }
	
	// Print the modified model that was returned from the execution.
    fmt.Println(model)
}
```

## Helper functions for strategies
The lib provides some helper functions to make your life easier, you may want to run
strategies conditionally for example.

### Group
You can use the `Group` helper function to group multiple strategies into a single strategy.

```go
func UsernameCorrect(ctx context.Context, container Container, model DataModel) (bool, error) {
    return model.ExternalData.UserName == "John Doe", nil
}

plan := speedrail.Plan(
    speedrail.If(
        UsernameCorrect, // Condition
        speedrail.Group( // Group several strategies into one strategy, this will be executed if the condition is not met
            SetUserName, // Strategy
            InsertUserToDatabase, // Strategy
        )
    )
)
```

### If
You can use the `If` helper function to run a strategy if a condition is met.

```go
func UsernameCorrect(ctx context.Context, container Container, model DataModel) (bool, error) {
    return model.ExternalData.UserName == "John Doe", nil
}

plan := speedrail.Plan(
    speedrail.If(
        UsernameCorrect, // Condition
        InsertUserToDatabase, // Strategy
    )
)
```

### IfElse
You can use the `IfElse` helper function to run a strategy if a condition is met, otherwise run another condition.

```go
func UsernameCorrect(ctx context.Context, container Container, model DataModel) (bool, error) {
    return model.ExternalData.UserName == "John Doe", nil
}

plan := speedrail.Plan(
    speedrail.IfElse(
        UsernameCorrect, // Condition
        InsertUserToDatabase, // Strategy if condition is met
        speedrail.Group( // Group several strategies into one strategy, this will be executed if the condition is not met
            SetUserName, // Strategy
            InsertUserToDatabase, // Strategy
        )
    )
)
```

### Merge
You can use the `Merge` helper function to run multiple strategies without breaking on error, instead it merges the
errors and return all of them together. It can be useful on features such as data validation.

```go
plan := speedrail.Plan(
    speedrail.Merge(
        InsertUserToDatabase, // Strategy if condition is met
        AnotherStrategyWithError // Strategy returns an error
    )
)
```

### ThrowError
You can use the `ThrowError` helper function to throw an error and stop the execution of the plan.

```go
plan := speedrail.Plan(
    speedrail.If(
        UserNameDoesNotExist, // Strategy if condition is met
        speedrail.ThrowError[any](speedrail.NewError(errors.New("user name does not exist"), http.StatusBadRequest, "user name does not exist")), // Strategy returns an error
    )
)
```

## Conditions

### Condition signature
Conditions are more limited in scope than a strategy, and used specifically to put conditions on if a strategy should be
executed or not. A condition can only be based on data present in the data model. A condition must follow the
`speedrail.Condition[M any]` signature, where `M` stands for `Model`. This Model must be the same as the model you used
in your strategies.

```go
type ConditionSignature = speedrail.Condition[Model]
```

### Example condition

```go
func condition1(model Model) bool {
    return model.ExternalData.UserName == "John Doe"
}

// Type checking so that condition follows our defined signature
var _ ConditionSignature = condition1
```

### Helper functions for conditions

### And
You can use the `And` helper function to run multiple conditions and return true if all of them are met.

```go
plan := speedrail.Plan(
    speedrail.If(
        speedrail.And(
            condition1, // Must be met
            speedrail.Or( // Either of condition2 or condition3 has to be met.
                condition2, // Condition
                condition3, // Condition
            ),
        ), // Condition argument
        DoSomethingStrategy, // Strategy if condition is met
    )
)
```

### Or
You can use the `Or` helper function to run multiple conditions and return true if any of them are met.
```go
plan := speedrail.Plan(
    speedrail.If(
        speedrail.Or(
            condition1, // Either condition1 is met
            speedrail.And( // Or both of condition2 and condition3 are met.
                condition2, // Condition
                condition3, // Condition
            ),
        ), // Condition argument
        DoSomethingStrategy, // Strategy if condition is met
    )
)
```

### Not
Not will invert the result of a condition.

```go
plan := speedrail.Plan(
    speedrail.If(
		speedrail.Not(condition1), // If condition is not met, succeed.
        DoSomethingStrategy, // Strategy if condition is met
    )
)
```