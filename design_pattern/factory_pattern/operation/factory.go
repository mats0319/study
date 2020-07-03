package mario

type OperationFactory interface {
    CreateOperation() Operation
}
