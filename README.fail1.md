Command

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

```
[flink-runner-job-invoker-2] ERROR org.apache.beam.runners.jobsubmission.JobInvocation - Error during job invocation go0job0101655507201275497696-red-0617230641-8e320fc2_ac2c4bdc-598d-4571-9392-c740b54aaf82.
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
	at java.base/jdk.internal.reflect.NativeMethodAccessorImpl.invoke0(Native Method)
	at java.base/jdk.internal.reflect.NativeMethodAccessorImpl.invoke(NativeMethodAccessorImpl.java:62)
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
Caused by: java.lang.RuntimeException: Emitting the record caused an I/O exception: Failed to serialize element. Serialized size (> 650016384 bytes) exceeds JVM heap space
	at org.apache.flink.runtime.operators.shipping.OutputCollector.collect(OutputCollector.java:69)
	at org.apache.flink.runtime.operators.util.metrics.CountingCollector.collect(CountingCollector.java:35)
	at org.apache.beam.runners.flink.translation.functions.HashingFlinkCombineRunner.combine(HashingFlinkCombineRunner.java:117)
	at org.apache.beam.runners.flink.translation.functions.FlinkReduceFunction.reduce(FlinkReduceFunction.java:108)
	at org.apache.flink.api.java.operators.translation.PlanUnwrappingReduceGroupOperator$TupleUnwrappingNonCombinableGroupReducer.reduce(PlanUnwrappingReduceGroupOperator.java:120)
	at org.apache.flink.runtime.operators.GroupReduceDriver.run(GroupReduceDriver.java:145)
	at org.apache.flink.runtime.operators.BatchTask.run(BatchTask.java:519)
	at org.apache.flink.runtime.operators.BatchTask.invoke(BatchTask.java:360)
	at org.apache.flink.runtime.taskmanager.Task.runWithSystemExitMonitoring(Task.java:958)
	at org.apache.flink.runtime.taskmanager.Task.restoreAndInvoke(Task.java:937)
	at org.apache.flink.runtime.taskmanager.Task.doRun(Task.java:766)
	at org.apache.flink.runtime.taskmanager.Task.run(Task.java:575)
	at java.base/java.lang.Thread.run(Thread.java:829)
Caused by: java.io.IOException: Failed to serialize element. Serialized size (> 650016384 bytes) exceeds JVM heap space
	at org.apache.flink.core.memory.DataOutputSerializer.resize(DataOutputSerializer.java:312)
	at org.apache.flink.core.memory.DataOutputSerializer.write(DataOutputSerializer.java:124)
	at org.apache.beam.runners.flink.translation.wrappers.DataOutputViewWrapper.write(DataOutputViewWrapper.java:44)
	at java.base/java.io.DataOutputStream.write(DataOutputStream.java:107)
	at java.base/java.io.FilterOutputStream.write(FilterOutputStream.java:108)
	at org.apache.beam.sdk.coders.ByteArrayCoder.encode(ByteArrayCoder.java:67)
	at org.apache.beam.sdk.coders.ByteArrayCoder.encode(ByteArrayCoder.java:56)
	at org.apache.beam.sdk.coders.ByteArrayCoder.encode(ByteArrayCoder.java:41)
	at org.apache.beam.sdk.coders.IterableLikeCoder.encode(IterableLikeCoder.java:113)
	at org.apache.beam.sdk.coders.IterableLikeCoder.encode(IterableLikeCoder.java:59)
	at org.apache.beam.sdk.coders.Coder.encode(Coder.java:136)
	at org.apache.beam.sdk.coders.KvCoder.encode(KvCoder.java:73)
	at org.apache.beam.sdk.coders.KvCoder.encode(KvCoder.java:37)
	at org.apache.beam.sdk.util.WindowedValue$FullWindowedValueCoder.encode(WindowedValue.java:607)
	at org.apache.beam.sdk.util.WindowedValue$FullWindowedValueCoder.encode(WindowedValue.java:598)
	at org.apache.beam.sdk.util.WindowedValue$FullWindowedValueCoder.encode(WindowedValue.java:558)
	at org.apache.beam.runners.flink.translation.types.CoderTypeSerializer.serialize(CoderTypeSerializer.java:110)
	at org.apache.flink.runtime.plugable.SerializationDelegate.write(SerializationDelegate.java:54)
	at org.apache.flink.runtime.io.network.api.writer.RecordWriter.serializeRecord(RecordWriter.java:132)
	at org.apache.flink.runtime.io.network.api.writer.RecordWriter.emit(RecordWriter.java:106)
	at org.apache.flink.runtime.io.network.api.writer.ChannelSelectorRecordWriter.emit(ChannelSelectorRecordWriter.java:54)
	at org.apache.flink.runtime.operators.shipping.OutputCollector.collect(OutputCollector.java:65)
	... 12 more
```