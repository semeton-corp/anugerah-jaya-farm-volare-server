package enum

type ExpenseCategory uint8

const (
	ExpenseCategoryUnknown     ExpenseCategory = 0
	ExpenseCategoryOperational ExpenseCategory = 1
	ExpenseCategoryOther       ExpenseCategory = 2
)

var (
	ExpenseCategoryMap = map[ExpenseCategory]string{
		ExpenseCategoryOperational: "Operasional",
		ExpenseCategoryOther:       "Lain-lain",
	}
)

func (c ExpenseCategory) String() string {
	return ExpenseCategoryMap[c]
}

func ValueOfExpenseCategory(value string) ExpenseCategory {
	for k, v := range ExpenseCategoryMap {
		if v == value {
			return k
		}
	}
	return ExpenseCategoryUnknown
}

func (c ExpenseCategory) IsValid() bool {
	_, ok := ExpenseCategoryMap[c]
	return ok
}
