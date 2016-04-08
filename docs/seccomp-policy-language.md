# Seccomp definition policy language

## Top level syntax

Each line is its own unit of parsing - there exists no way of extending expressions over multiple lines.

Every line can be one of several types - specifically, they can be assignments, rules or comments.

## Comments

A comment will start with a literal octothorpe (#) in column 0, and continues until the end of the line
No processing of comments will happen.

## Assignments

(TODO)

## Rules

A rule can take several different forms. Each rule will be for one specific systemcall. That systemcall will be referred to by its common name. There can only be one rule per systemcall for each policy file. A rule can result in either a boolean result, or a direct return action.
If a boolean result happens, the rule will generate a return action based on the DEFAULT positive or negative action for that policy file. Specifically, a positive result from the rule, will return the DEFAULT POSITIVE action, and the negative result will return the DEFAULT NEGATIVE action. 

There are several different format for rules. They all start with the name of the system call, followed by possible spaces, followed by a colon and possible spaces. The first form takes a boolean expression after the colon, and will generate a positive or negative result depending on the outcome of that boolean expression:

  read: arg0 == 1 || 42 + 5 == arg1

The second form allows you to return a specific error number when a system call is invoked:

  read: return 42

The third form combines these two, in such a way that if the given expression is positive, the default positive action will be returned, but if negative, the error number specified by return will be used:

  read: arg0==1; return 55

(TODO, continue here)

## Syntax of numbers

Numbers can be represented in four different formats, following the standard conventions:

Binary: 0b010101010
Octal: 0777
Decimal: 42
Hexadecimal: 0xFEFE or 0XfeFE

All numbers represent 32 bit unsigned numbers. Negative numbers can only be represented implicitly, through arithmetic operations.

Comparisons with 64bit arguments will be handled automatically when comparing against literal numbers. 
If it's necessary to compare the upper 32 bits of an argument, explicit bit masking and shifting needs to be done.
For example:
  (arg0 >> 32) == 0b01001

## Syntax of expressions

In all expressions, there will be variables named arg0 to arg5 available. These are 64 bit unsigned numbers. All comparisons with these numbers will only operate on the lower 32 bits, such that the upper 32bits have to be 0 for the comparison to succeed.

The arguments to unary or binary operators can be any VALUE, where VALUE is defined to either be one of the argument names, an explicit number, or another expression.

### Arithmetic

All standard arithmetic operators from C are available and follow the same precedence rules. Specifically, these operators are:
- Parenthesis
- Plus (+)
- Minus (-)
- Multiplication (*)
- Division (/)
- Binary and (&)
- Binary or (|)
- Binary xor (^)
- Binary negation (~)
- Left shift (<<)
- Right shift (>>)
- Modulo (%)

### Boolean operations

The outcome of every rule will be defined by boolean operations. These primarily include comparisons of various kinds. Boolean operations support these operators:
- Parenthesis
- Boolean OR (||)
- Boolean AND(&&)
- Boolean negation (!)
- Comparison operators
  - Equal (==)
  - Not equal (!=)
  - Greater than (>)
  - Greater or equal to (>=)
  - Less than (<)
  - Less than or equal to (<=)
  - Bits set (this operator will mask the left hand side with the right hand, and return true if the result has any bits set) (&)
- Inclusion:
  arg0 in [1,2,3,4]
  arg0 not in [1, 2, 3, 4]
  There are no general arrays in the language - square brackets can only be used in conjunction with the in and not in operator
  the in/not in operators are not case sensitive. Any valid value or name can be used inside the square brackets. Values have to be separated
  by commas, and arbitrary amount of whitespace (tabs or spaces). The in/not in operator is the only binary operator where the two sides are not equivalent.

THese can all be arbitrarily nested.

In a boolean context, the strings "1" and "true" can be used as literal true values, while "0" and "false" can be used as literal false values. This is mostly useful to create short rules like:
  read: 1
  fcntl: false

