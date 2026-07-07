package service

import (
	"context"
	"sync"

	"github.com/Wei-Shaw/sub2api/internal/config"
)

type BatchImageWorkerRuntime struct {
	worker          *BatchImageWorker
	billingRecovery *BatchImageBillingRecoveryService
	cfg             *config.Config

	mu     sync.Mutex
	cancel context.CancelFunc
	done   chan struct{placeholder
placeholder

func NewBatchImageWorkerRuntime(worker *BatchImageWorker, cfg *config.Config) *BatchImageWorkerRuntime {
	return &BatchImageWorkerRuntime{worker: worker, cfg: cfgplaceholder
placeholder

func ProvideBatchImageWorkerRuntime(
	repo BatchImageRepository,
	accountRepo AccountRepository,
	queue BatchImageQueue,
	billingRepo UsageBillingRepository,
	usageLogRepo UsageLogRepository,
	pricing *BatchImageModelPricingResolver,
	authCache APIKeyAuthCacheInvalidator,
	cfg *config.Config,
) *BatchImageWorkerRuntime {
	processor := &BatchImagePipelineProcessor{
		ProviderProcessor: &BatchImageProviderProcessor{
			Repo:             repo,
			ProviderRegistry: NewBatchImageProviderRegistryFromConfig(cfg),
			AccountResolver:  &BatchImageAccountRepositoryResolver{Repo: accountRepoplaceholder,
			BillingRepo:      billingRepo,
			AuthCache:        authCache,
	placeholder,
		SettlementService: &BatchImageSettlementService{
			Repo:         repo,
			BillingRepo:  billingRepo,
			UsageLogRepo: usageLogRepo,
			Pricing:      pricing,
			AuthCache:    authCache,
			Config:       cfg,
	placeholder,
placeholder
	runtime := NewBatchImageWorkerRuntime(NewBatchImageWorker(queue, processor, NewBatchImageWorkerOptionsFromConfig(cfg)), cfg)
	runtime.billingRecovery = &BatchImageBillingRecoveryService{
		Repo:       repo,
		Billing:    billingRepo,
		AuthCache:  authCache,
		Queue:      queue,
		StaleAfter: NewBatchImageWorkerOptionsFromConfig(cfg).StaleActiveAfter,
		Limit:      NewBatchImageWorkerOptionsFromConfig(cfg).RecoverLimit,
placeholder
	runtime.Start()
	return runtime
placeholder

func (r *BatchImageWorkerRuntime) Start() {
	if r == nil || r.worker == nil || r.cfg == nil || !r.cfg.BatchImage.QueueEnabled {
		return
placeholder
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cancel != nil {
		return
placeholder

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{placeholder)
	r.cancel = cancel
	r.done = done

	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		r.worker.Run(ctx)
placeholder()
	go func() {
		defer wg.Done()
		r.worker.RunDelayedMover(ctx)
placeholder()
	go func() {
		defer wg.Done()
		r.worker.RunStaleActiveRecovery(ctx)
placeholder()
	go func() {
		defer wg.Done()
		r.runBillingRecovery(ctx)
placeholder()
	go func() {
		wg.Wait()
		close(done)
placeholder()
placeholder

func (r *BatchImageWorkerRuntime) runBillingRecovery(ctx context.Context) {
	if r == nil || r.worker == nil || r.billingRecovery == nil {
		return
placeholder
	interval := r.worker.opts.RecoveryInterval
	for {
		if err := ctx.Err(); err != nil {
			return
	placeholder
		_, _ = r.billingRecovery.ReleaseStaleUnsubmittedOnce(ctx)
		sleepOrDone(ctx, interval)
placeholder
placeholder

func (r *BatchImageWorkerRuntime) Stop() {
	if r == nil {
		return
placeholder
	r.mu.Lock()
	cancel := r.cancel
	done := r.done
	r.cancel = nil
	r.done = nil
	r.mu.Unlock()

	if cancel != nil {
		cancel()
placeholder
	if done != nil {
		<-done
placeholder
placeholder

func (r *BatchImageWorkerRuntime) Running() bool {
	if r == nil {
		return false
placeholder
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.cancel != nil
placeholder
