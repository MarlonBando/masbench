Comparison
==========

This guide explains how to compare benchmark results using masbench's comparison tools to analyze the performance differences between different algorithm implementations or configurations.

Overview
--------

The masbench comparison feature allows you to:

- Compare two benchmark results side by side
- Generate visual charts for different performance metrics
- Create detailed difference reports
- Identify performance patterns and improvements

Basic Comparison
----------------

To compare two benchmark results, use the ``masbench compare`` command:

.. code-block:: bash

   masbench compare benchmark1-name benchmark2-name

For example, to compare an A* implementation against a BFS implementation:

.. code-block:: bash

   masbench compare astar-v1 bfs-v1

.. important::
   Both benchmark results must exist in your benchmark folder. If either benchmark is not found, masbench will display an error message.

Prerequisites
-------------

Before running comparisons, ensure you have:

1. At least two completed benchmark runs
2. Valid CSV result files for both benchmarks
3. The benchmark names match the folder names in your benchmark directory

Generated Output
----------------

When you run a comparison, masbench creates a new folder structure:

.. code-block:: text

   benchmark-results/
   └── comparisons/
       └── benchmark1vsbenchmark2/
           ├── Generated.png
           ├── Explored.png
           ├── MemoryAlloc.png
           ├── Time.png
           ├── Actions.png
           └── benchmark1vsbenchmark2_table.png

Output Files Explained
~~~~~~~~~~~~~~~~~~~~~~

**Metric Charts** (``*.png``)
   Individual bar charts comparing each performance metric:

   - ``Generated.png``: Nodes generated during search
   - ``Explored.png``: Nodes explored during search  
   - ``MemoryAlloc.png``: Memory allocated during execution
   - ``Time.png``: Execution time for each level
   - ``Actions.png``: Number of actions in the solution

**Difference Report** (``*_table.png``)
   A comprehensive table showing side-by-side comparison of all metrics for each level, including difference calculations.

Difference Report Table
~~~~~~~~~~~~~~~~~~~~~~~

The difference report provides a detailed breakdown showing how the benchmark1 compared to benchmark2.
If a metric is colored in green, it indicates benchmark1 outperformed benchmark2 for that metric on that level.

.. image:: ../build/html/_static/table_comparison_example.png
   :alt: Example of benchmark comparison table
   :align: center
   :width: 100%

Next Steps
----------

After analyzing your comparisons:

1. Identify the best-performing algorithm for your use case
2. Look for opportunities to combine strengths from different approaches
3. Plan further optimizations based on the insights gained
4. Document your findings for future reference

.. seealso::
   - For running benchmarks, see the :doc:`running_benchmarks` guide
   - For initial setup, see the :doc:`getting_started` guide
