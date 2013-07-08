# Saturate - to utilize all the resources

Saturate performs a multi-tier fanout of tasks to workers with an
indirection to separate the control between global concurrency and
per-worker type concurrency.

For example, given a collection of objects that can be retrieved from
any one of several servers, you can ensure that retrieval stays busy
by avoiding busier workers in favor of available workers.
