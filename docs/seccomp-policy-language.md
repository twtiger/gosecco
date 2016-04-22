# Seccomp definition policy language

## Top level syntax

Each line is its own unit of parsing - there exists no way of extending expressions over multiple lines.

Every line can be one of several types - specifically, they can be assignments, rules or comments.

In general, each line will be parsed and understood in the context of only the previous lines. That means that variables and macros have to be defined before used. This also stops recursive actions from being possible.

## Valid names

In order to simplify implementation, we reuse the parser from the Go programming language. That means that certain words will not be valid as variable names. These include the typical Golang keywords such as "for", "type", "if", "func" and so on.

## Comments

A comment will start with a literal octothorpe (#) in column 0, and continues until the end of the line
No processing of comments will happen.

## Default actions

Each rule can generate a positive or a negative action, depending on whether the boolean result of that rule is positive or negative. When compiling the program it is possible to set the defaults that should be used. This might not always be the most convenient option though, so the language also supports defining default actions inside of the file itself. These can be specified by assigning the special values DEFAULT_POSITIVE and DEFAULT_NEGATIVE in the usual manner of assignment. The standard actions available have mnemonic names as well. These are  "trap", "kill", "allow", "trace". If a number is given, this will be interpreted as returning an ERRNO action for that number:

  DEFAULT_POSITIVE = trace
  DEFAULT_NEGATIVE = 42

It is suggested to define these at the top of the file to minimize confusion. It is theoretically possible to change the default actions through the file, but that is discouraged.

## Assignments
  
Assignments allow the policy writer to simplify and extract complex arithmetic operations. The operational semantics of the assignment is as if the expression had been put inline at the place where the variable is referenced. The expression defining the variable has to be well formed in isolation, but can refer to previously defined variables. The compiler will perform arithmetic simplification on all expressions in order to reduce the number of operations needed at runtime.

  var1 = 412 * 3

## Macros

If a variable expression refers to any of the special argument variables, or contains any boolean operators or the return operator, the assignment will instead refer to a macro. The operational semantics for a macro is the same as for variables. Any expression that refers to a macro becomes a macro.

Macros can take arguments - the argument list follows the usual rules and the evaluation will use simple alpha renaming before compilation.

  var2 = arg0 == 5 && arg1 == 42; return 6
  f(x) = x == 5
  g(y) = y == 6

  read: f(arg0) || g(arg0) || g(arg1) 
  read: var2

## Rules

A rule can take several different forms. Each rule will be for one specific systemcall. That systemcall will be referred to by its common name. There can only be one rule per systemcall for each policy file. A rule can result in either a boolean result, or a direct return action.
If a boolean result happens, the rule will generate a return action based on the DEFAULT positive or negative action for that policy file. Specifically, a positive result from the rule, will return the DEFAULT POSITIVE action, and the negative result will return the DEFAULT NEGATIVE action.

There are several different format for rules. They all start with the name of the system call, followed by possible spaces, followed by a colon and possible spaces. The first form takes a boolean expression after the colon, and will generate a positive or negative result depending on the outcome of that boolean expression:

  read: arg0 == 1 || 42 + 5 == arg1

The second form allows you to return a specific error number when a system call is invoked:

  read: return 42

The third form combines these two, in such a way that if the given expression is positive, the default positive action will be returned, but if negative, the error number specified by return will be used:

  read: arg0==1; return 55

Rules can specify their own custom positive and negative actions that differ from the default. This uses the same naming convention as the default actions described above. The syntax for describing them is simple:

  read[+trace, -kill] : 1 == 2
  read[+42] : arg0 == 1
  read[-55] : arg0 > 1
  
The order of the actions is arbitrary, and either part can be left out. The plus sign signifies the positive action, and the minus the negative action. If no actions are specified, the square brackets can be left off, and the default actions for the file will be used.

## Syntax of numbers

Numbers can be represented in four different formats, following the standard conventions:

Octal: 0777
Decimal: 42
Hexadecimal: 0xFEFE or 0XfeFE

All numbers represent 32 bit unsigned numbers. Negative numbers can only be represented implicitly, through arithmetic operations.

Comparisons with 64bit arguments will be handled automatically when comparing against literal numbers.
If it's necessary to compare the upper 32 bits of an argument, explicit bit masking and shifting needs to be done.
For example:
  (arg0 >> 32) == 0xfe

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
- Binary negation (^)
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
  in(arg0, 1,2,3,4)
  notIn(arg0, 1, 2, 3, 4)
  the in/not in operators are not case sensitive. Any valid value or name can be used inside the brackets. Values have to be separated
  by commas, and arbitrary amount of whitespace (tabs or spaces). The in/notIn operator is the function like application that is not actually a function

These can all be arbitrarily nested.

In a boolean context, the strings "1" and "true" can be used as literal true values, while "0" and "false" can be used as literal false values. This is mostly useful to create short rules like:
  read: 1
  fcntl: false
