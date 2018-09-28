Language Reference
==================

Comments
--------

**Single Line Comment**

.. image:: http://ferdinandsilva.com/static/comment1.png

**Multiline Comment**

.. image:: http://ferdinandsilva.com/static/comment2.png

Variables
---------

**Rules:**

- A variable name must start with a letter or the underscore character
- A variable name cannot start with a number
- A variable name can only contain alpha-numeric characters and underscores (A-z, 0-9, !, and _ )
- Variable names are case-sensitive (name and nAme are two different variables)

**Data Types:**

- *Integer* - a non-decimal number between âˆ’2,147,483,647 and 2,147,483,647

    .. image:: http://ferdinandsilva.com/static/integer.png

- *String* - a sequence of characters. A string can be any text inside quotes. You can use single or double quotes.

    .. image:: http://ferdinandsilva.com/static/string.png

- *Float* - a number with a decimal point.

    .. image:: http://ferdinandsilva.com/static/float.png

- *Nil* - a special data type that represents a variable with no value

    .. image:: http://ferdinandsilva.com/static/nil.png

- *Boolean* - a value that is either TRUE or FALSE

    .. image:: http://ferdinandsilva.com/static/boolean.png

- *Lineup* - is an ordered sequence of items. All the items in a lineup do not need to be of the same type.

    .. image:: http://ferdinandsilva.com/static/lineup.png

- *Glossary* - is an unordered collection of items. It has a key: value pair.

    .. image:: http://ferdinandsilva.com/static/glossary.png

**Constant:**

It has a name starting with an uppercase character.

.. image:: http://ferdinandsilva.com/static/constant.png

Operators
---------

**Arithmetic Operators**

+----------+-----------------+-----------------------------------------------------------------------+
| Operator |       Name      |                  Example                                              |
+==========+=================+=======================================================================+
|      \+  |    Addition     | .. image:: http://ferdinandsilva.com/static/sum.png                   |
+----------+-----------------+-----------------------------------------------------------------------+
|      \-  |   Subtraction   | .. image:: http://ferdinandsilva.com/static/difference.png            |
+----------+-----------------+-----------------------------------------------------------------------+
|      \*  | Multiplication  | .. image:: http://ferdinandsilva.com/static/product.png               |
+----------+-----------------+-----------------------------------------------------------------------+
|       /  |   Division      | .. image:: http://ferdinandsilva.com/static/quotient.png              |
+----------+-----------------+-----------------------------------------------------------------------+

**Comparison Operators**

+----------+-----------------------------+-----------------------------------------------------------------------+
| Operator |       Name                  |                  Example                                              |
+==========+=============================+=======================================================================+
|    \=\=  | Equal                       | .. image:: http://ferdinandsilva.com/static/equal.png                 |
+----------+-----------------------------+-----------------------------------------------------------------------+
|     >    | Greater than                | .. image:: http://ferdinandsilva.com/static/greater.png               |
+----------+-----------------------------+-----------------------------------------------------------------------+
|     <    | Less than                   | .. image:: http://ferdinandsilva.com/static/less.png                  |
+----------+-----------------------------+-----------------------------------------------------------------------+
|     <\=  | Less than or Equal          | .. image:: http://ferdinandsilva.com/static/less_or_equal.png         |
+----------+-----------------------------+-----------------------------------------------------------------------+
|     >\=  | Greater than or Equal       | .. image:: http://ferdinandsilva.com/static/g_or_equal.png            |
+----------+-----------------------------+-----------------------------------------------------------------------+

If Statement
------------

**If Statement**

Executes some code if one condition is true.

**Syntax**
::
    if (condition)
        code to be executed if condition is true
    fi

**Example**

.. image:: http://ferdinandsilva.com/static/equal.png

**If...Else Statement**

Executes some code if condition is true and another code if that condition is false.

**Syntax**
::
    if (condition)
        code to be executed if condition is true
    el 
        code to be executed if condition is false
    fi

**Example**

.. image:: http://ferdinandsilva.com/static/ifelse2.png

**If...ElseIf...Else Statement**

Executes different codes for more than two conditions.

**Syntax**
::
    if (condition)
        code to be executed if this condition is true
    ef (condition)
        code to be executed if this condition is true
    el
        code to be executed if all conditions are false
    fi

**Example**

.. image:: http://ferdinandsilva.com/static/ifelseif2.png

Looping Statements
------------------

**For Loop**

Execute a block of code a specified number of times where start counter is lower than end counter.

**Syntax**
::
    fl (start counter, end counter)
        code to be executed
    lf

**Example**

.. image:: http://ferdinandsilva.com/static/forloop.png

**Break Statement**

When a break statement is encountered inside a loop, the loop is immediately terminated and the program control resumes at the next statement following the loop.

**Example**

.. image:: http://ferdinandsilva.com/static/break2.png

Functions
---------

**Syntax**
::
    fd functionName(parameter1, parameter2)
        code to be executed
    df

**Example 1 (Function without parameter & return)**

.. image:: http://ferdinandsilva.com/static/function1.png

**Example 2 (Function with parameters & return)**

.. image:: http://ferdinandsilva.com/static/function2.png

Built-in Functions
------------------

Console Functions
~~~~~~~~~~~~~~~~~

- **p** - writes string to the standard output (stdout).

    **Declaration:**
    ::
        string p(string)

- **r** - presents a prompt to the user and read a string from standard input (stdin).

    **Declaration:**
    ::
        string r(string)

- **sc** - a function to set the text color on a console screen.

    **Declaration:**
    ::
        Nil sc(integer 0..8)

Conversion Functions
~~~~~~~~~~~~~~~~~~~~

- **toi** - convert any data type to integer type.

    **Declaration:**
    ::
        integer toi(any data type)

- **tos** - convert any data type to string type.

    **Declaration:**
    ::
        string tos(any data type)

String Functions
~~~~~~~~~~~~~~~~

- **str_rpl** - returns a copy of the first parameter in which the occurrences of second parameter have been replaced with third parameter.

    **Declaration:**
    ::
        string str_rpl(string, string, string)

- **str_spl** - split a string into a lineup.

    **Declaration:**
    ::
        lineup str_spl(string, string)

System Functions
~~~~~~~~~~~~~~~~

- **ex** - terminates program execution and returns the status value to the system.

    **Declaration:**
    ::
        Nil ex(integer)

- **abt** - print a message and exit the current script.

    **Declaration:**
    ::
        Nil abt(string)

- **exe** - executes an internal operating system command.

    **Declaration:**
    ::
        glossary exe(string)

- **zzz** - delays program execution for a given number of milliseconds.

    **Declaration:**
    ::
        Nil zzz(integer)

- **sav** - returns raw command-line arguments.

    **Declaration:**
    ::
        lineup sav()

Other Functions
~~~~~~~~~~~~~~~

- **!** - (reverse boolean), converts boolean data type true to false and vice versa.

    **Declaration:**
    ::
        bool !(bool)

- **len** - returns the item count of a glossary or lineup variable.

    **Declaration:**
    ::
        integer len(any)

- **i** - used to include a titik file in another file.

    **Declaration:**
    ::
        Nil i(string)