# gcsim AST documentation

## Variables

Variable types are inferred from the first usage

## Built in function

## ExprStmt

These are basic statements to be executed, they can be broken down into:

- Actions
- Variable Assignment
- Math expressions
- Conditional expressions
- Function call

Expressions all return a value

- Actions return true or false representing if it was executed
- Variable assignment returns the value that was assigned
- Math expressions returns the results of the expression
- Conditional expressions return the results of the conditions
- Function call returns the signature of the function

```go
type ExprStmt interface {
    Eval() interface{}
}
```

### Condition

This is a basic conditional statement

```go
type CondStmt struct {
	Left   *CondStmt
	Right  *CondStmt
	IsLeaf bool
	Op     string //&& || ( )
	Expr   Condition
}
```

## WhileStmt

This is a basic loop

```
            while
           /     \
          /       \
     condition     \
                 []expr
```

## Functions

Functions are just a map of function name ot a tree I guess?
