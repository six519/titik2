Language Reference
==================

Hello World Code
----------------

::

    ^
        (Hello World Code)
        Source code at: https://github.com/six519/titik2
    ^

    strVariable = "Mabuhay Ka"

    #function definition
    fd hello_world(var)
        intVariable = 0
        #loop statement
        fl (0, 8)
            intVariable = intVariable + 1
            sc(intVariable) #change text color
            p(var) #It will print "Mabuhay Ka" 5 times
            zzz(2000) #sleep for 2 seconds
        lf
        
        sc(0) #reset the text color
    df

    hello_world(strVariable) #call function

Hello World Code Output
-----------------------

.. image:: http://ferdinandsilva.com/static/titik_output.png

Comments
--------

**Single Line Comment**

.. image:: http://ferdinandsilva.com/static/comment1.png

**Multiline Comment**

.. image:: http://ferdinandsilva.com/static/comment3.png

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

**Assignment Operator**

+----------+----------------------------------------------------------------------------------------+
| Operator |                  Example                                                               |
+==========+========================================================================================+
|      \=  | .. image:: http://ferdinandsilva.com/static/assignment.png                             |
+----------+----------------------------------------------------------------------------------------+

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
|    \=\=  | Equal                       | .. image:: http://ferdinandsilva.com/static/equal2.png                |
+----------+-----------------------------+-----------------------------------------------------------------------+
|     <>   | Not Equal                   | .. image:: http://ferdinandsilva.com/static/not_equal.png             |
+----------+-----------------------------+-----------------------------------------------------------------------+
|     >    | Greater than                | .. image:: http://ferdinandsilva.com/static/greater.png               |
+----------+-----------------------------+-----------------------------------------------------------------------+
|     <    | Less than                   | .. image:: http://ferdinandsilva.com/static/less.png                  |
+----------+-----------------------------+-----------------------------------------------------------------------+
|     <\=  | Less than or Equal          | .. image:: http://ferdinandsilva.com/static/less_or_equal.png         |
+----------+-----------------------------+-----------------------------------------------------------------------+
|     >\=  | Greater than or Equal       | .. image:: http://ferdinandsilva.com/static/g_or_equal.png            |
+----------+-----------------------------+-----------------------------------------------------------------------+

**Logical Operators**

+----------+-----------------------------+-----------------------------------------------------------------------+
| Operator |       Name                  |                  Example                                              |
+==========+=============================+=======================================================================+
|    \|    | OR                          | .. image:: http://ferdinandsilva.com/static/or.png                    |
+----------+-----------------------------+-----------------------------------------------------------------------+
|     &    | AND                         | .. image:: http://ferdinandsilva.com/static/and.png                   |
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

.. image:: http://ferdinandsilva.com/static/equal2.png

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

**Example 1**

.. image:: http://ferdinandsilva.com/static/forloop2.png

**Example 2 (Infinite Loop)**

.. image:: http://ferdinandsilva.com/static/infinite.png

**While Loop**

Execute a statement or code block repeatedly as long as an expression is true.

**Syntax**
::
    wl (condition)
        code to be executed
    lw

**Example**

.. image:: http://ferdinandsilva.com/static/while.png

**For Each Loop**

Used to iterate through elements of lineup.

**Syntax**
::
    fea (lineup variable, variable to hold the item)
        code to be executed
    aef

**Example**

.. image:: http://ferdinandsilva.com/static/foreach.png

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

- **rp** - prompt the user to input character/text without echoing.

    **Declaration:**
    ::
        string rp(string)

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

- **str_l** - converts a string to lowercase letters.

    **Declaration:**
    ::
        string str_l(string)

- **str_u** - converts a string to uppercase letters.

    **Declaration:**
    ::
        string str_u(string)

- **str_t** - removes whitespace from the left and right side of a string.

    **Declaration:**
    ::
        string str_t(string)

- **str_chr** - returns the character that represents the specified code point.

    **Declaration:**
    ::
        string str_chr(integer)

- **str_ord** - returns the code point of a specified character.

    **Declaration:**
    ::
        integer str_ord(string)

- **str_sub** - extracts parts of a string.

    **Declaration:**
    ::
        string str_sub(string, integer, integer)

- **str_ind** - get index of specified substring.

    **Declaration:**
    ::
        integer str_ind(string, string)

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

- **gcp** - get current working directory.

    **Declaration:**
    ::
        string gcp()

File Processing Functions
~~~~~~~~~~~~~~~~~~~~~~~~~

- **fo** - open file for reading/writing/appending.

    **Declaration:**
    ::
        string fo(string, string)

- **fc** - close file.

    **Declaration:**
    ::
        Nil fc(string)

- **fw** - write string file.

    **Declaration:**
    ::
        Nil fw(string, string)

- **fr** - read the file using specific bytes count.

    **Declaration:**
    ::
        string fr(string, integer)

- **flrm** - removes path and any children it contains.

    **Declaration:**
    ::
        bool flrm(string)

- **flmv** - moves old path to new path.

    **Declaration:**
    ::
        bool flmv(string, string)

- **flcp** - copy old path to new path.

    **Declaration:**
    ::
        bool flcp(string, string)

Other Functions
~~~~~~~~~~~~~~~

- **!** - (reverse boolean), converts boolean data type true to false and vice versa.

    **Declaration:**
    ::
        bool !(bool)

- **len** - returns the item count of a glossary/lineup/string variable.

    **Declaration:**
    ::
        integer len(any)

- **i** - used to include a titik file in another file.

    **Declaration:**
    ::
        Nil i(string)

- **in** - used to check if a variable is a Nil type.

    **Declaration:**
    ::
        bool in(any)

- **la** - add new item to a lineup.

    **Declaration:**
    ::
        lineup la(lineup, any)

- **lp** - pop item from a slice through its index.

    **Declaration:**
    ::
        lineup lp(lineup, integer)

- **sgv** - set global variable.

    **Declaration:**
    ::
        Nil sgv(string, any)

MySQL Functions
~~~~~~~~~~~~~~~

- **mysql_set** - initialize MySQL connection and returns the connection reference string.

    **Declaration:**
    ::
        string mysql_set(string, string, string, string)

- **mysql_q** - executes SQL statement.

    **Declaration:**
    ::
        bool mysql_q(string, string)

- **mysql_cr** - cleanup MySQL resources.

    **Declaration:**
    ::
        Nil mysql_cr(string)

- **mysql_fa** - get a result row as a lineup by column name.

    **Declaration:**
    ::
        lineup mysql_fa(string, string)

SQLite Functions
~~~~~~~~~~~~~~~

- **sqlite_set** - set SQLite file to open.

    **Declaration:**
    ::
        Nil sqlite_set(string)

- **sqlite_q** - executes SQL statement.

    **Declaration:**
    ::
        bool sqlite_q(string)

- **sqlite_cr** - cleanup SQLite resources.

    **Declaration:**
    ::
        Nil sqlite_cr()

- **sqlite_fa** - get a result row as a lineup by column name.

    **Declaration:**
    ::
        lineup sqlite_fa(string)

HTTP Functions
~~~~~~~~~~~~~~

- **http_au** - registers the Titik function for the given URL pattern.

    **Declaration:**
    ::
        Nil http_au(string, string)

- **http_su** - set the static directory for the given URL pattern.

    **Declaration:**
    ::
        Nil http_su(string, string)

- **http_run** - starts an HTTP server with a given address.

    **Declaration:**
    ::
        Nil http_run(string)

- **http_p** - print a string to a web browser.

    **Declaration:**
    ::
        Nil http_p(string)

- **http_gq** - parses query string and returns the corresponding values.

    **Declaration:**
    ::
        lineup http_gq(string)

- **http_gfp** - returns HTTP POST parameter.

    **Declaration:**
    ::
        lineup http_gfp(string)

- **http_lt** - loads HTML template file.

    **Declaration:**
    ::
        Nil http_lt(string, glossary)

- **http_gp** - get current URL path.

    **Declaration:**
    ::
        string http_gp()

- **http_h** - set HTTP header.

    **Declaration:**
    ::
        Nil http_h(string, string)

- **http_cr** - HTTP client.

    **Declaration:**
    ::
        glossary http_cr(string, string, glossary, glossary)

Cryptographic Functions
~~~~~~~~~~~~~~~~~~~~~~~

- **m5** - get MD5 sum of a string.

    **Declaration:**
    ::
        string m5(string)

- **s1** - get SHA1 sum of a string.

    **Declaration:**
    ::
        string s1(string)

- **s256** - get SHA256 sum of a string.

    **Declaration:**
    ::
        string s256(string)

- **s512** - get SHA512 sum of a string.

    **Declaration:**
    ::
        string s512(string)

- **b64e** - Encode the string using Base64.

    **Declaration:**
    ::
        string b64e(string)

- **b64d** - Decode the Base64 string.

    **Declaration:**
    ::
        string b64d(string)

- **rsae** - Encrypt the string by RSA public key.

    **Declaration:**
    ::
        string rsae(string, string)

Socket Functions
~~~~~~~~~~~~~~~~

- **netc** - initiates TCP/UDP client connection and returns the connection reference string.

    **Declaration:**
    ::
        string netc(string, string)

- **netl** - initiates TCP server connection and returns the connection reference string.

    **Declaration:**
    ::
        string netl(string, string)

- **netul** - initiates UDP server connection and returns the connection reference string.

    **Declaration:**
    ::
        string netul(string, string, integer)

- **netulf** - initiates UDP server connection and call a function handler.

    **Declaration:**
    ::
        string netulf(string, string, integer, string)

- **netla** - accepts TCP client connection and returns the connection reference string.

    **Declaration:**
    ::
        string netla(string)

- **netlaf** - accepts TCP client connection and call a function handler.

    **Declaration:**
    ::
        string netlaf(string, string)

- **netlx** - closes server TCP socket connection.

    **Declaration:**
    ::
        Nil netlx(string)

- **netx** - closes client socket connection.

    **Declaration:**
    ::
        Nil netx(string)

- **netw** - transmits TCP/UDP message.

    **Declaration:**
    ::
        Nil netw(string, string)

- **netr** - receives TCP/UDP message.

    **Declaration:**
    ::
        string netr(string, integer)

- **netur** - receives UDP message.

    **Declaration:**
    ::
        string netur(string, integer)

Hello World Code (Web)
----------------------

::

    fd index()
        http_p("<h1>Hello World</h1>")
    df

    http_au("/", "index")
    http_run(":8080")

Hello World Code (Web) Output
-----------------------------

.. image:: http://ferdinandsilva.com/static/web.png