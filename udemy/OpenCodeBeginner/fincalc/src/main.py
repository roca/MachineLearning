"""
Compound Interest Calculator CLI

An interactive command-line tool that calculates compound interest
with regular monthly contributions. Prompts the user for investment
details and outputs the results as formatted JSON.

Outputs:
    - Final balance after the investment period
    - Total contributions made
    - Total interest earned
"""

import json
import sys
from typing import TypedDict


class CompoundInterestResult(TypedDict):
    """Type definition for the calculation results."""

    final_balance: float
    total_contributions: float
    total_interest_earned: float


def get_float_input(prompt: str, allow_zero: bool = True) -> float:
    """
    Prompt the user for a float value with input validation.

    Re-prompts until the user enters a valid non-negative number.

    Args:
        prompt: The message to display to the user.
        allow_zero: Whether to allow zero as a valid input.

    Returns:
        A valid float value entered by the user.
    """
    while True:
        try:
            value = float(input(prompt))
            if not allow_zero and value == 0:
                print("Error: Value must be greater than zero.")
                continue
            if value < 0:
                print("Error: Value cannot be negative.")
                continue
            return value
        except ValueError:
            print("Error: Please enter a valid number.")


def get_int_input(prompt: str, min_val: int = 1) -> int:
    """
    Prompt the user for an integer value with input validation.

    Re-prompts until the user enters a valid whole number.

    Args:
        prompt: The message to display to the user.
        min_val: The minimum allowed value.

    Returns:
        A valid integer value entered by the user.
    """
    while True:
        try:
            value = int(input(prompt))
            if value < min_val:
                print(f"Error: Value must be at least {min_val}.")
                continue
            return value
        except ValueError:
            print("Error: Please enter a valid whole number.")


def get_compounding_frequency() -> str:
    """
    Prompt the user to select a compounding frequency from a menu.

    Displays numbered options for daily, monthly, or yearly compounding
    and validates the user's selection.

    Returns:
        The selected compounding frequency as a string ('daily', 'monthly', or 'yearly').
    """
    print("\nCompounding Frequency Options:")
    print("  1. Daily (365 times per year)")
    print("  2. Monthly (12 times per year)")
    print("  3. Yearly (1 time per year)")

    frequency_map = {"1": "daily", "2": "monthly", "3": "yearly"}

    while True:
        choice = input("\nEnter your choice (1/2/3): ").strip()
        if choice in frequency_map:
            return frequency_map[choice]
        print("Error: Please enter 1, 2, or 3.")


def compound_interest_calulator() -> None:
    """
    Interactive compound interest calculator.

    Prompts the user for investment details via the CLI:
    1. Principal amount - the initial investment
    2. Monthly contributions - recurring deposits added each month
    3. Annual interest rate - the yearly interest rate as a percentage
    4. Compounding frequency - how often interest is applied (daily, monthly, yearly)
    5. Time period - the investment duration in years

    Then calculates and displays the final balance, total contributions,
    and total interest earned as formatted JSON.

    The calculation uses month-by-month simulation for accuracy, applying
    an effective monthly rate derived from the nominal annual rate and
    the chosen compounding frequency.
    """
    print("=" * 50)
    print("    COMPOUND INTEREST CALCULATOR")
    print("=" * 50)

    # --- Gather user inputs via interactive prompts ---
    print("\nEnter the following investment details:\n")

    principal = get_float_input("  Principal amount ($): ", allow_zero=True)
    monthly_contribution = get_float_input(
        "  Monthly contribution ($): ", allow_zero=True
    )
    annual_rate = get_float_input(
        "  Annual interest rate (e.g., 7.5 for 7.5%): ", allow_zero=True
    )
    compounding_frequency = get_compounding_frequency()
    time_years = get_int_input("\n  Time period (years): ", min_val=1)

    # --- Perform the compound interest calculation ---
    # Convert annual rate from percentage to decimal
    rate_decimal = annual_rate / 100

    # Map compounding frequency to the number of compounding periods per year
    periods_per_year = {"daily": 365, "monthly": 12, "yearly": 1}
    n = periods_per_year[compounding_frequency]

    # Total number of months in the investment period
    total_months = time_years * 12

    # Simulate month-by-month to accurately handle monthly contributions
    balance = principal
    total_contributions = 0.0

    for _ in range(total_months):
        # Add the monthly contribution to the balance
        balance += monthly_contribution
        total_contributions += monthly_contribution

        # Calculate the effective monthly interest rate from the nominal annual rate
        # This accounts for the compounding frequency within each month
        monthly_rate = (1 + rate_decimal / n) ** (n / 12) - 1

        # Apply interest for this month
        balance *= 1 + monthly_rate

    # Calculate total interest earned
    total_interest = balance - principal - total_contributions

    # --- Build the result dictionary ---
    result: CompoundInterestResult = {
        "final_balance": round(balance, 2),
        "total_contributions": round(total_contributions, 2),
        "total_interest_earned": round(total_interest, 2),
    }

    # --- Display results as JSON ---
    print("\n" + "=" * 50)
    print("    RESULTS")
    print("=" * 50)
    print(json.dumps(result, indent=2))
    print("=" * 50)


if __name__ == "__main__":
    try:
        compound_interest_calulator()
    except KeyboardInterrupt:
        print("\n\nCalculator cancelled by user.")
        sys.exit(0)
    except EOFError:
        print("\n\nNo input received. Exiting.")
        sys.exit(1)
