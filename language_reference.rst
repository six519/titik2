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
- A variable name can only contain alpha-numeric characters and underscores (A-z, 0-9, and _ )
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

**Constant:**

It has a name starting with an uppercase character.

.. image:: http://ferdinandsilva.com/static/constant.png

Operators
---------

**Arithmetic Operators**

+----------+-----------------+-----------------------------------------------------------------------+
| Operator |       Name      |                  Example                                              |
+==========+=================+=======================================================================+
|      \+  |    Addition     | .. image:: http://ferdinandsilva.com/static/sum.png                    |
+----------+-----------------+-----------------------------------------------------------------------+
|      \-  |   Subtraction   | .. image:: http://ferdinandsilva.com/static/difference.png             |
+----------+-----------------+-----------------------------------------------------------------------+
|      \*  | Multiplication  | .. image:: http://ferdinandsilva.com/static/product.png                |
+----------+-----------------+-----------------------------------------------------------------------+
|       /  |   Division      | .. image:: http://ferdinandsilva.com/static/quotient.png               |
+----------+-----------------+-----------------------------------------------------------------------+

**Comparison Operators**

+----------+-----------------+-----------------------------------------------------------------------+
| Operator |       Name      |                  Example                                              |
+==========+=================+=======================================================================+
|    \=    | Equal           | .. image:: http://ferdinandsilva.com/static/equal.png                  |
+----------+-----------------+-----------------------------------------------------------------------+
|     >    | Greater than    | .. image:: http://ferdinandsilva.com/static/greater.png                |
+----------+-----------------+-----------------------------------------------------------------------+
|     <    | Less than       | .. image:: http://ferdinandsilva.com/static/less.png                   |
+----------+-----------------+-----------------------------------------------------------------------+

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
    e 
        code to be executed if condition is false
    fi

**Example**

.. image:: http://ferdinandsilva.com/static/ifelse.png

**If...ElseIf...Else Statement**

Executes different codes for more than two conditions.

**Syntax**
::
    if (condition)
        code to be executed if this condition is true
    ef (condition)
        code to be executed if this condition is true
    e
        code to be executed if all conditions are false
    fi

**Example**

.. image:: http://ferdinandsilva.com/static/ifelseif.png

Looping Statements
------------------

**For Loop**

Execute a block of code a specified number of times where start counter is lower than end counter.

**Syntax**
::
    fl (start counter to end counter)
        code to be executed
    lf

**Example**

.. image:: http://ferdinandsilva.com/static/forward.png

**Break Statement**

When a break statement is encountered inside a loop, the loop is imstatictely terminated and the program control resumes at the next statement following the loop.

**Example**

.. image:: http://ferdinandsilva.com/static/break.png

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

- **zzz** - delays program execution for a given number of milliseconds.

    **Declaration:**
    ::
        Nil zzz(integer)

- **p** - writes string to the standard output (stdout).

    **Declaration:**
    ::
        Nil p(string)

- **i** - used to include a titik file in another file.

    **Declaration:**
    ::
        Nil i(string)

- **tof** - convert string/integer to float type.

    **Declaration:**
    ::
        float tof(string or integer)

- **toi** - convert float/string to integer type.

    **Declaration:**
    ::
        integer toi(string or float)

- **tos** - convert float/integer to string type.

    **Declaration:**
    ::
        string tos(float or integer)

- **ex** - terminates program execution and returns the status value to the system.

    **Declaration:**
    ::
        Nil ex(integer)

- **sc** - a function to set the text color on a console screen.

    **Declaration:**
    ::
        Nil sc(integer 0..7)

- **flcp** - makes a copy of the file source to destination. If successfull, the return is the destination.

    **Declaration:**
    ::
        string flcp(string, string)

- **flmv** - moves the file source to destination. If successfull, the return is the destination.

    **Declaration:**
    ::
        string flmv(string, string)

- **flrm** - deletes a file. If successfull, the return is 1, if not then it will return 0.

    **Declaration:**
    ::
        integer flrm(string)

- **exe** - executes an internal operating system command. If successfull, the return is 1, if not then it will return 0.

    **Declaration:**
    ::
        integer exe(string)

- **r** - presents a prompt to the user and read a string from standard input (stdin).

    **Declaration:**
    ::
        string r(string)

- **rnd** - return a random integer between 0 and a specified max number.

    **Declaration:**
    ::
        integer rnd(integer)

- **sac** - return count of command line arguments.

    **Declaration:**
    ::
        integer sac()

- **savf** - return the first command line argument passed to a Titik script.

    **Declaration:**
    ::
        string savf()

- **rndstr** - return a random string with a length specified by a parameter.

    **Declaration:**
    ::
        string rndstr(integer)