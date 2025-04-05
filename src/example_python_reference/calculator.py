"""
simple_math.calculator

A simple calculator module for basic arithmetic operations.
"""


def add(a: float, b: float) -> float:
    """
    Add two numbers.

    Parameters
    ----------
    a : float
        The first number.
    b : float
        The second number.

    Returns
    -------
    float
        The sum of a and b.

    Examples
    --------
    >>> add(1, 2)
    3.0
    """
    return float(a + b)


def subtract(a: float, b: float) -> float:
    """
    Subtract one number from another.

    Parameters
    ----------
    a : float
        The number to subtract from.
    b : float
        The number to subtract.

    Returns
    -------
    float
        The result of a - b.

    Examples
    --------
    >>> subtract(5, 3)
    2.0
    """
    return float(a - b)


def multiply(a: float, b: float) -> float:
    """
    Multiply two numbers.

    Parameters
    ----------
    a : float
        The first number.
    b : float
        The second number.

    Returns
    -------
    float
        The product of a and b.

    Examples
    --------
    >>> multiply(2, 4)
    8.0
    """
    return float(a * b)


def divide(a: float, b: float) -> float:
    """
    Divide one number by another.

    Parameters
    ----------
    a : float
        The numerator.
    b : float
        The denominator.

    Returns
    -------
    float
        The result of a / b.

    Raises
    ------
    ZeroDivisionError
        If b is zero.

    Examples
    --------
    >>> divide(10, 2)
    5.0
    """
    if b == 0:
        raise ZeroDivisionError("Division by zero is not allowed.")
    return float(a / b)
