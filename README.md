## Reproduction of OOM issue

This is a reproduction of a beam pipeline crashing due to OOM.
https://github.com/apache/beam/issues/21817

Start up the Flink server

```shell
bazel run :flink-job-server-driver -- \
  --flink-conf-dir $PWD/flink-conf/ \
  --job-port 8099
```

```shell
mkdir -p tmp
```

### Successful cases

```shell
bazel run -c opt //cmd/example_pipeline --   --runner flink   --environment_type LOOPBACK   --endpoint localhost:8099   --output "$PWD/tmp/records"   --record-bytes 10000   --record-count 87600 --alsologtostderr
```

### Problem 1: 1,000 x 5MB records

1,000 x 5MB records fails per [README.fail1.md](README.fail1.md)

```shell
bazel run //cmd/example_pipeline -- \
  --alsologtostderr \
  --runner flink \
  --environment_type LOOPBACK \
  --endpoint localhost:8099 \
  --output "$PWD/tmp/records" \
  --record-bytes 5000000 \
  --record-count 10000
```


### Problem 2:

The problem I am trying to reproduce in this project is this error, caused by
another pipeline: `java.io.IOException: Cannot write record to fresh sort
buffer. Record too large.`

```
org.apache.flink.runtime.client.JobExecutionException: Job execution failed.
	at org.apache.flink.runtime.jobmaster.JobResult.toJobExecutionResult(JobResult.java:144)
	at org.apache.flink.runtime.minicluster.MiniClusterJobClient.lambda$getJobExecutionResult$3(MiniClusterJobClient.java:137)
	at java.base/java.util.concurrent.CompletableFuture$UniApply.tryFire(CompletableFuture.java:642)
	at java.base/java.util.concurrent.CompletableFuture.postComplete(CompletableFuture.java:506)
	at java.base/java.util.concurrent.CompletableFuture.complete(CompletableFuture.java:2073)
	at org.apache.flink.runtime.rpc.akka.AkkaInvocationHandler.lambda$invokeRpc$1(AkkaInvocationHandler.java:258)
	at java.base/java.util.concurrent.CompletableFuture.uniWhenComplete(CompletableFuture.java:859)
	at java.base/java.util.concurrent.CompletableFuture$UniWhenComplete.tryFire(CompletableFuture.java:837)
	at java.base/java.util.concurrent.CompletableFuture.postComplete(CompletableFuture.java:506)
	at java.base/java.util.concurrent.CompletableFuture.complete(CompletableFuture.java:2073)
	at org.apache.flink.util.concurrent.FutureUtils.doForward(FutureUtils.java:1389)
	at org.apache.flink.runtime.concurrent.akka.ClassLoadingUtils.lambda$null$1(ClassLoadingUtils.java:93)
	at org.apache.flink.runtime.concurrent.akka.ClassLoadingUtils.runWithContextClassLoader(ClassLoadingUtils.java:68)
	at org.apache.flink.runtime.concurrent.akka.ClassLoadingUtils.lambda$guardCompletionWithContextClassLoader$2(ClassLoadingUtils.java:92)
	at java.base/java.util.concurrent.CompletableFuture.uniWhenComplete(CompletableFuture.java:859)
	at java.base/java.util.concurrent.CompletableFuture$UniWhenComplete.tryFire(CompletableFuture.java:837)
	at java.base/java.util.concurrent.CompletableFuture.postComplete(CompletableFuture.java:506)
	at java.base/java.util.concurrent.CompletableFuture.complete(CompletableFuture.java:2073)
	at org.apache.flink.runtime.concurrent.akka.AkkaFutureUtils$1.onComplete(AkkaFutureUtils.java:47)
	at akka.dispatch.OnComplete.internal(Future.scala:300)
	at akka.dispatch.OnComplete.internal(Future.scala:297)
	at akka.dispatch.japi$CallbackBridge.apply(Future.scala:224)
	at akka.dispatch.japi$CallbackBridge.apply(Future.scala:221)
	at scala.concurrent.impl.CallbackRunnable.run(Promise.scala:60)
	at org.apache.flink.runtime.concurrent.akka.AkkaFutureUtils$DirectExecutionContext.execute(AkkaFutureUtils.java:65)
	at scala.concurrent.impl.CallbackRunnable.executeWithValue(Promise.scala:68)
	at scala.concurrent.impl.Promise$DefaultPromise.$anonfun$tryComplete$1(Promise.scala:284)
	at scala.concurrent.impl.Promise$DefaultPromise.$anonfun$tryComplete$1$adapted(Promise.scala:284)
	at scala.concurrent.impl.Promise$DefaultPromise.tryComplete(Promise.scala:284)
	at akka.pattern.PromiseActorRef.$bang(AskSupport.scala:621)
	at akka.pattern.PipeToSupport$PipeableFuture$$anonfun$pipeTo$1.applyOrElse(PipeToSupport.scala:24)
	at akka.pattern.PipeToSupport$PipeableFuture$$anonfun$pipeTo$1.applyOrElse(PipeToSupport.scala:23)
	at scala.concurrent.Future.$anonfun$andThen$1(Future.scala:532)
	at scala.concurrent.impl.Promise.liftedTree1$1(Promise.scala:29)
	at scala.concurrent.impl.Promise.$anonfun$transform$1(Promise.scala:29)
	at scala.concurrent.impl.CallbackRunnable.run(Promise.scala:60)
	at akka.dispatch.BatchingExecutor$AbstractBatch.processBatch(BatchingExecutor.scala:63)
	at akka.dispatch.BatchingExecutor$BlockableBatch.$anonfun$run$1(BatchingExecutor.scala:100)
	at scala.runtime.java8.JFunction0$mcV$sp.apply(JFunction0$mcV$sp.java:12)
	at scala.concurrent.BlockContext$.withBlockContext(BlockContext.scala:81)
	at akka.dispatch.BatchingExecutor$BlockableBatch.run(BatchingExecutor.scala:100)
	at akka.dispatch.TaskInvocation.run(AbstractDispatcher.scala:49)
	at akka.dispatch.ForkJoinExecutorConfigurator$AkkaForkJoinTask.exec(ForkJoinExecutorConfigurator.scala:48)
	at java.base/java.util.concurrent.ForkJoinTask.doExec(ForkJoinTask.java:290)
	at java.base/java.util.concurrent.ForkJoinPool$WorkQueue.topLevelExec(ForkJoinPool.java:1020)
	at java.base/java.util.concurrent.ForkJoinPool.scan(ForkJoinPool.java:1656)
	at java.base/java.util.concurrent.ForkJoinPool.runWorker(ForkJoinPool.java:1594)
	at java.base/java.util.concurrent.ForkJoinWorkerThread.run(ForkJoinWorkerThread.java:183)
Caused by: org.apache.flink.runtime.JobException: Recovery is suppressed by NoRestartBackoffTimeStrategy
	at org.apache.flink.runtime.executiongraph.failover.flip1.ExecutionFailureHandler.handleFailure(ExecutionFailureHandler.java:138)
	at org.apache.flink.runtime.executiongraph.failover.flip1.ExecutionFailureHandler.getFailureHandlingResult(ExecutionFailureHandler.java:82)
	at org.apache.flink.runtime.scheduler.DefaultScheduler.handleTaskFailure(DefaultScheduler.java:252)
	at org.apache.flink.runtime.scheduler.DefaultScheduler.maybeHandleTaskFailure(DefaultScheduler.java:242)
	at org.apache.flink.runtime.scheduler.DefaultScheduler.updateTaskExecutionStateInternal(DefaultScheduler.java:233)
	at org.apache.flink.runtime.scheduler.SchedulerBase.updateTaskExecutionState(SchedulerBase.java:684)
	at org.apache.flink.runtime.scheduler.SchedulerNG.updateTaskExecutionState(SchedulerNG.java:79)
	at org.apache.flink.runtime.jobmaster.JobMaster.updateTaskExecutionState(JobMaster.java:444)
	at jdk.internal.reflect.GeneratedMethodAccessor50.invoke(Unknown Source)
	at java.base/jdk.internal.reflect.DelegatingMethodAccessorImpl.invoke(DelegatingMethodAccessorImpl.java:43)
	at java.base/java.lang.reflect.Method.invoke(Method.java:566)
	at org.apache.flink.runtime.rpc.akka.AkkaRpcActor.lambda$handleRpcInvocation$1(AkkaRpcActor.java:316)
	at org.apache.flink.runtime.concurrent.akka.ClassLoadingUtils.runWithContextClassLoader(ClassLoadingUtils.java:83)
	at org.apache.flink.runtime.rpc.akka.AkkaRpcActor.handleRpcInvocation(AkkaRpcActor.java:314)
	at org.apache.flink.runtime.rpc.akka.AkkaRpcActor.handleRpcMessage(AkkaRpcActor.java:217)
	at org.apache.flink.runtime.rpc.akka.FencedAkkaRpcActor.handleRpcMessage(FencedAkkaRpcActor.java:78)
	at org.apache.flink.runtime.rpc.akka.AkkaRpcActor.handleMessage(AkkaRpcActor.java:163)
	at akka.japi.pf.UnitCaseStatement.apply(CaseStatements.scala:24)
	at akka.japi.pf.UnitCaseStatement.apply(CaseStatements.scala:20)
	at scala.PartialFunction.applyOrElse(PartialFunction.scala:123)
	at scala.PartialFunction.applyOrElse$(PartialFunction.scala:122)
	at akka.japi.pf.UnitCaseStatement.applyOrElse(CaseStatements.scala:20)
	at scala.PartialFunction$OrElse.applyOrElse(PartialFunction.scala:171)
	at scala.PartialFunction$OrElse.applyOrElse(PartialFunction.scala:172)
	at scala.PartialFunction$OrElse.applyOrElse(PartialFunction.scala:172)
	at akka.actor.Actor.aroundReceive(Actor.scala:537)
	at akka.actor.Actor.aroundReceive$(Actor.scala:535)
	at akka.actor.AbstractActor.aroundReceive(AbstractActor.scala:220)
	at akka.actor.ActorCell.receiveMessage(ActorCell.scala:580)
	at akka.actor.ActorCell.invoke(ActorCell.scala:548)
	at akka.dispatch.Mailbox.processMailbox(Mailbox.scala:270)
	at akka.dispatch.Mailbox.run(Mailbox.scala:231)
	at akka.dispatch.Mailbox.exec(Mailbox.scala:243)
	... 5 more
Caused by: java.io.IOException: Cannot write record to fresh sort buffer. Record too large.
	at org.apache.flink.runtime.operators.chaining.SynchronousChainedCombineDriver.collect(SynchronousChainedCombineDriver.java:190)
	at org.apache.flink.runtime.operators.util.metrics.CountingCollector.collect(CountingCollector.java:35)
	at org.apache.flink.runtime.operators.chaining.ChainedMapDriver.collect(ChainedMapDriver.java:78)
	at org.apache.flink.runtime.operators.util.metrics.CountingCollector.collect(CountingCollector.java:35)
	at org.apache.beam.runners.flink.translation.functions.FlinkExecutableStagePruningFunction.flatMap(FlinkExecutableStagePruningFunction.java:58)
	at org.apache.beam.runners.flink.translation.functions.FlinkExecutableStagePruningFunction.flatMap(FlinkExecutableStagePruningFunction.java:30)
	at org.apache.flink.runtime.operators.FlatMapDriver.run(FlatMapDriver.java:113)
	at org.apache.flink.runtime.operators.BatchTask.run(BatchTask.java:519)
	at org.apache.flink.runtime.operators.BatchTask.invoke(BatchTask.java:360)
	at org.apache.flink.runtime.taskmanager.Task.runWithSystemExitMonitoring(Task.java:958)
	at org.apache.flink.runtime.taskmanager.Task.restoreAndInvoke(Task.java:937)
	at org.apache.flink.runtime.taskmanager.Task.doRun(Task.java:766)
	at org.apache.flink.runtime.taskmanager.Task.run(Task.java:575)
	at java.base/java.lang.Thread.run(Thread.java:829)
```