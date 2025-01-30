import pandas as pd
import matplotlib.pyplot as plt

csv_file = "fermentation_results.csv"
df = pd.read_csv(csv_file)

# Clean up column names
df.columns = df.columns.str.strip()

# Convert first column to numeric
time_col = df.columns[0]
df[time_col] = pd.to_numeric(df[time_col], errors='coerce')

# Plot each measurement separately in its own window
for col in df.columns[1:]:
    plt.figure(figsize=(8, 5))  # Create a new figure for each plot
    plt.plot(df[col], df[time_col], marker='o', linestyle='-', linewidth=0.8)
    plt.xlabel(col)
    plt.ylabel("Time (s)")
    plt.title(f"Fermentation Process: {col} vs Time")
    plt.grid()
    plt.show(block=False)  # Show the plot without blocking execution

input("Press Enter to close all plots...")
plt.close('all')  # Close all windows when the user presses Enter
