import pandas as pd
import matplotlib.pyplot as plt

# Load CSV file
csv_file = "fermentation_results.csv"
df = pd.read_csv(csv_file)

# Clean up column names
df.columns = df.columns.str.strip()

# Convert time column to numeric
time_col = df.columns[0]
df[time_col] = pd.to_numeric(df[time_col], errors='coerce')

# Restrict time range to between 1000s and 5000s
df = df[(df[time_col] >= 1000) & (df[time_col] <= 5000)]

# Identify the pH column (if available)
ph_col = [col for col in df.columns if "pH" in col or col.lower() == "ph"]
ph_col = ph_col[0] if ph_col else None  # Take the first matching pH column if found

# Plot each measurement separately with respect to time
for col in df.columns[1:]:  # Skip the first column (time)
    plt.figure(figsize=(8, 5))  # Create a new figure for each plot
    plt.plot(df[time_col], df[col], marker='o', linestyle='-', linewidth=1.2)

    plt.xlabel("Time (s)")
    plt.ylabel(col)
    plt.title(f"Fermentation Process: {col} vs Time")
    plt.grid()

    # Adjust pH plot y-axis range
    if col == ph_col:
        plt.ylim(4.0, 6.7)

    plt.show(block=False)  # Show the plot without blocking execution

# Wait for user input before closing all plots
input("Press Enter to close all plots...")
plt.close('all')  # Close all windows when the user presses Enter
