Running Benchmarks
==================

This guide explains how to run a benchmark and understand the results.

Prerequisites
-------------

Before running benchmarks, ensure you have:

1. Completed the initialization process (see :doc:`getting_started`)
2. A valid ``masbench_config.yml`` file in your project root
3. Your server executable (``.jar`` file) accessible
4. Test levels in your specified levels directory
5. A working client implementation

Basic Benchmark Execution
-------------------------

To run a benchmark, use the ``masbench run`` command followed by a benchmark name:

.. code-block:: bash

   masbench run my-first-benchmark

The benchmark name is used to:

- Create a unique folder for the results
- Name the output files
- Organize multiple benchmark runs

.. important::
   Benchmark names must be unique. If you try to run a benchmark with a name that already exists, masbench will display an error and exit.

Adding Notes to Your Benchmark
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

You can add a descriptive message to your benchmark run using the ``-m`` or ``--message`` flag:

.. code-block:: bash

   masbench run algorithm-v2 -m "Testing A* with improved heuristic"

This message helps you remember what changes you were testing when reviewing results later.

Output Structure
----------------

After running a benchmark, you'll find the following structure in your benchmark folder:

.. code-block:: text

   benchmarks/
   └── my-first-benchmark/
       ├── logs/
       │   ├── my-first-benchmark_server.zip
       │   └── my-first-benchmark_client.clog
       └── my-first-benchmark_results.csv

File Descriptions
~~~~~~~~~~~~~~~~~

**Server Logs** (``*_server.zip``)
   Contains detailed server execution logs, including level loading, client communication, and any server-side errors.

**Client Logs** (``*_client.clog``)
   Raw output from your client, including debug information, algorithm progress, and any client-side errors.

**Results CSV** (``*_results.csv``)
   Processed benchmark data in CSV format with the following columns:

   - ``LevelName``: Name of the level file
   - ``Solved``: Whether the level was solved (true/false)
   - ``Actions``: Number of actions in the solution
   - ``Time``: Execution time in milliseconds
   - ``Generated``: Number of nodes generated during search
   - ``Explored``: Number of nodes explored during search
   - ``MemoryAlloc``: Memory allocated during execution
   - ``MaxAlloc``: Peak memory allocation

Default Output vs Extended Metrics
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

By default, masbench only captures the following basic metrics from your benchmark runs:

- ``LevelName``: Name of the level file  
- ``Solved``: Whether the level was solved (true/false)
- ``Actions``: Number of actions in the solution
- ``Time``: Execution time in milliseconds

The additional metrics (``Generated``, ``Explored``, ``MemoryAlloc``, and ``MaxAlloc``) are **not included in the output by default**. 

To track these extended metrics, your client must output specific strings at the end of execution as comments. Each string must be on a new line:

.. code-block:: text

   #Explored:
   #Generated:
   #Alloc:

When masbench detects these strings in your client output, it will parse the values and include them in the results CSV. If these strings are not present, the corresponding columns will be empty or contain default values.

.. note::
   The output must be formatted as comments (starting with ``#``). For example, in Python you should add something like:
   
   .. code-block:: python
   
      print("#Explored: 123", flush=True)
      print("#Generated: 1039", flush=True)
      print("#Alloc: 200", flush=True)
   
   This will produce output that masbench can parse:
   
   .. code-block:: text
   
      [client][message] #Explored: 156
      [client][message] #Generated: 234
      [client][message] #Alloc: 3072

Example Results
~~~~~~~~~~~~~~~

Here's an example of what the CSV results might look like:

.. code-block:: text

   LevelName,Solved,Actions,Time,Generated,Explored,MemoryAlloc,MaxAlloc
   SAsoko1_01.lvl,true,12,45,127,89,2048,4096
   SAsoko1_02.lvl,true,18,78,234,156,3072,6144
   SAsoko1_03.lvl,false,0,300000,5670,4321,8192,16384

Troubleshooting Common Issues
-----------------------------

Benchmark Already Exists
~~~~~~~~~~~~~~~~~~~~~~~~

If you see this error:

.. code-block:: text

   Error: Benchmark with name 'my-benchmark' already exists. Please remove it before running a new one.

You have several options:

1. **Choose a different name**: Use a new benchmark name
2. **Remove the existing benchmark**: Delete the folder manually
3. **Rename the existing benchmark**: Move it to preserve the results

.. code-block:: bash

   # Remove existing benchmark
   rm -rf benchmark-results/my-benchmark
   
   # Or rename it to preserve results
   mv benchmark-results/my-benchmark benchmark-results/my-benchmark-old

Server Not Found
~~~~~~~~~~~~~~~~

If you get an error about the server not being found:

1. Check that the ``ServerPath`` in your config points to the correct file
2. Ensure the ``.jar`` file exists and is accessible
3. Verify you have Java installed and available in your PATH

Client Command Issues
~~~~~~~~~~~~~~~~~~~~

If your client fails to run:

1. Test your client command manually first
2. Ensure all dependencies are installed
3. Check that the client command in your config is correct
4. Review the client logs for specific error messages

Performance Tips
----------------

Selecting Test Levels
~~~~~~~~~~~~~~~~~~~~~

- Start with a small set of representative levels
- Include levels of varying difficulty
- Consider creating a separate "quick test" folder for rapid iteration
- Use the full level set for final benchmarks

Next Steps
----------

Once you have benchmark results:

1. Analyze the CSV data to identify performance patterns
2. Compare different algorithm implementations
3. Use the comparison tools to visualize differences

.. seealso::
   - For comparing benchmark results, see the :doc:`comparison` guide
   - For initial setup, see the :doc:`getting_started` guide
