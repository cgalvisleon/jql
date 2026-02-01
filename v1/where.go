package jql

import "github.com/cgalvisleon/jql/jql"

/**
* Eq
* @param field string, value interface{}
* @return jql.Condition
**/
func Eq(field string, value interface{}) *jql.Condition {
	return jql.Eq(field, value)
}

/**
* Neg
* @param field string, value interface{}
* @return jql.Condition
**/
func Neg(field string, value interface{}) *jql.Condition {
	return jql.Neg(field, value)
}

/**
* Less
* @param field string, value interface{}
* @return jql.Condition
**/
func Less(field string, value interface{}) *jql.Condition {
	return jql.Less(field, value)
}

/**
* LessEq
* @param field string, value interface{}
* @return jql.Condition
**/
func LessEq(field string, value interface{}) *jql.Condition {
	return jql.LessEq(field, value)
}

/**
* More
* @param field string, value interface{}
* @return jql.Condition
**/
func More(field string, value interface{}) *jql.Condition {
	return jql.More(field, value)
}

/**
* MoreEq
* @param field string, value interface{}
* @return jql.Condition
**/
func MoreEq(field string, value interface{}) *jql.Condition {
	return jql.MoreEq(field, value)
}

/**
* Like
* @param field string, value interface{}
* @return jql.Condition
**/
func Like(field string, value interface{}) *jql.Condition {
	return jql.Like(field, value)
}

/**
* In
* @param field string, value []interface{}
* @return jql.Condition
**/
func In(field string, value []interface{}) *jql.Condition {
	return jql.In(field, value)
}

/**
* NotIn
* @param field string, value []interface{}
* @return jql.Condition
**/
func NotIn(field string, value []interface{}) *jql.Condition {
	return jql.NotIn(field, value)
}

/**
* Is
* @param field string, value []interface{}
* @return jql.Condition
**/
func Is(field string, value []interface{}) *jql.Condition {
	return jql.Is(field, value)
}

/**
* IsNot
* @param field string, value []interface{}
* @return jql.Condition
**/
func IsNot(field string, value []interface{}) *jql.Condition {
	return jql.IsNot(field, value)
}

/**
* Null
* @param field string
* @return jql.Condition
**/
func Null(field string) *jql.Condition {
	return jql.Null(field)
}

/**
* NotNull
* @param field string
* @return jql.Condition
**/
func NotNull(field string) *jql.Condition {
	return jql.NotNull(field)
}

/**
* Between
* @param field string, min, max interface{}
* @return jql.Condition
**/
func Between(field string, min, max interface{}) *jql.Condition {
	return jql.Between(field, min, max)
}

/**
* NotBetween
* @param field string, min, max interface{}
* @return jql.Condition
**/
func NotBetween(field string, min, max interface{}) *jql.Condition {
	return jql.NotBetween(field, min, max)
}
