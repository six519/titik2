titik2
======

Titik Programming Language/Interpreter Written In Go Lang

Intro
=====

This is for experimental and educational purpose only. ;)

Compile From Source (Windows)
=============================
::
    
    go build -tags win .

Compile From Source (Other OS)
==============================
::
    
    go build .

Binary Release (Installation Instructions)
==========================================

Click here_

.. _here: https://github.com/six519/titik2/blob/master/install.rst

Hello World Code
================
::

    ^
        Multiline comment
        (Hello World Code)
        Source code at: https://github.com/six519/titik2
    ^

    floatVariable = 25.55
    strVariable = "Mabuhay " + 'Ka' #concatenation

    #function definition
    fd hello_world(var)
        intVariable = 0
        #if statement
        if(var == 'Mabuhay Ka')

            #loop statement
            fl (0, 8)
                intVariable = intVariable + 1
                sc(intVariable) #change text color
                p(var + '\n') #It will print "Mabuhay Ka" 8 times
                zzz(2000) #sleep for 2 seconds
            lf
            
            sc(0) #reset the text color
        el
            #else
            p('Not Mabuhay Ka')
        fi
    df

    hello_world(strVariable) #call function

Hello World Code Output
=======================

.. image:: http://ferdinandsilva.com/static/titik_output.png

Language Reference
==================

Click me_

.. _me: https://github.com/six519/titik2/blob/master/language_reference.rst