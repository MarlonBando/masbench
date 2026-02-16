Changes
=======

This page documents the changes made to the masbench project over time.

Version 1.3.0 
-------------

**New Features:**

* **run** - Added ``-a`` / ``--algorithm`` flag to specify the algorithm to use during benchmark runs. This allows you to dynamically select algorithms (e.g., bfs, dfs, greedy, astar) without modifying your configuration file.

**Bug Fix**

* Fixed the bug that when the `Time to solve` output from the server was above 1000.000 seconds the parsing was breaking showing only `1`  

Version 1.2.2
-------------

* Fixed bug that was breaking the log parsing if a white space was present in the path of a level
* Fixed version unknown bug

Version 1.2.1
-------------

* Improvments to the doc

Version 1.2.0
-------------

**Improvements:**

* **list** - Now displays benchmark descriptions alongside names. Use ``--name-only`` flag to show only names.
* **run** - Benchmark descriptions are now saved to ``.md`` files for future reference.

Version 1.1.0
-------------

**New Features:**

* **list** - View all benchmarks in your benchmark folder
* **rm** - Remove a benchmark and all related comparisons
* **summary** - Generate HTML summary reports for one or more benchmarks

**Improvements:**

* **compare** - Now generates interactive HTML reports instead of static PNG images. Reports include sortable tables, filterable data, dark mode toggle, and chart export options.

**Bug Fixes:**

* Fixed action count display bug where values of 1000 or more showed only the first digit (e.g., 1000 displayed as 1)

Version 1.0.0 (Initial Release)
--------------------------------

This is the initial release of masbench.

**Features:**

* **init** - Initialize masbench configuration in a repository with interactive setup
* **run** - Execute benchmarks with optional message notes and performance tracking
* **compare** - Compare two benchmark results with visual tables and graphs
* **version** - Display the current version of masbench
