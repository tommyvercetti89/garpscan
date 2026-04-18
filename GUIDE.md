# GarpScan Usage Guide 📖

Welcome to the beginner-friendly guide for GarpScan! If you are new to Go or the GarpScan SDK, this document will help you understand how everything fits together.

## 🏗️ Core Concepts (Simply Explained)

Think of GarpScan as a high-powered car:

1.  **The Engine (`garpscan.Engine`):** This is the car itself. It handles the speed (concurrency) and manages the trip.
2.  **The Plugins (`garpscan.Plugin`):** These are the tools you pack in the car (e.g., Maps, Spare Tire). In GarpScan, these are the modules that know *how* to scan something (like a Port Scanner).
3.  **The Targets:** These are your destinations (IP addresses or Websites).
4.  **The Result (`garpscan.Result`):** This is the trip report—what happened at each destination.

---

## ⚡ "Cheat Sheet" for Python Users

If you are coming from Python, here is how GarpScan methods map to familiar concepts:

| GarpScan (Go) | Python Analogy | What it does |
| :--- | :--- | :--- |
| `garpscan.NewEngine()` | `socket.socket()` | "Starts the motor." Creates the main scanner. |
| `garpscan.WithWorkers(n)` | `pool_size = n` | Sets how many scans happen at the same time. |
| `engine.AddPlugin(p)` | `list.append(plugin)` | Adds a specific scanning feature to the tool. |
| `engine.Scan(ctx, targets)`| `for t in targets: start(t)` | **The Big Button.** Starts the actual scanning process. |
| `reporter.ExportJSON()` | `json.dump(file)` | Saves the messy data into a clean JSON file. |

---

## 🛠️ The 4-Step Workflow

To use GarpScan in your own project, follow these four steps:

### Step 1: Initialize the Engine
Create the motor. You usually want to set the number of "workers" (simultaneous scans).
```go
engine := garpscan.NewEngine(garpscan.WithWorkers(50))
```

### Step 2: Add Your Plugins
Tell the engine what kind of scanning you want to do. You can add many plugins at once.
```go
// Add a port scanner plugin
engine.AddPlugin(portscanner.New([]int{80, 443}, 2*time.Second))
```

### Step 3: Run the Scan
Provide your targets and start the process. This returns a "Channel" (a stream of data) that you can listen to.
```go
targets := []string{"1.1.1.1", "google.com"}
resultsChan := engine.Scan(context.Background(), targets)
```

### Step 4: Handle the Results
Use the `reporter` package to save or display what the scanner finds.
```go
// This will print everything to your screen in JSON format
reporter.ExportJSON(os.Stdout, resultsChan)
```

---

## 🚀 Pro Tips for Success

1.  **Context is King:** Always use a `context.WithTimeout` if you want the scan to stop automatically after a certain amount of time.
2.  **Worker Balance:** If you set `WithWorkers` too high (e.g., 5000), your internet connection might get choked. Start with 50-100 and experiment.

# Run all tests
make test           # (Linux/macOS)
.\tasks.ps1 test    # (Windows)

# Tidy module dependencies
make tidy           # (Linux/macOS)
.\tasks.ps1 tidy    # (Windows)

# Run the basic example
make example        # (Linux/macOS)
.\tasks.ps1 example  # (Windows)

3.  **Check the Examples:** Look at the [examples/](examples/) folder for a fully working script you can copy and paste!

---

## ❓ Need More Detail?
Check out the technical [API Reference](API.md) for a full list of every single function and variable available.
