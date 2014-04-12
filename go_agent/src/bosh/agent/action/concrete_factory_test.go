package action_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh/agent/action"
	fakeas "bosh/agent/applier/applyspec/fakes"
	fakeappl "bosh/agent/applier/fakes"
	fakecomp "bosh/agent/compiler/fakes"
	boshdrain "bosh/agent/drain"
	faketask "bosh/agent/task/fakes"
	fakeblobstore "bosh/blobstore/fakes"
	fakeinfrastructure "bosh/infrastructure/fakes"
	fakejobsuper "bosh/jobsupervisor/fakes"
	boshlog "bosh/logger"
	fakenotif "bosh/notification/fakes"
	fakeplatform "bosh/platform/fakes"
	boshntp "bosh/platform/ntp"
	fakesettings "bosh/settings/fakes"
)

func init() {
	Describe("concreteFactory", func() {
		var (
			settings            *fakesettings.FakeSettingsService
			platform            *fakeplatform.FakePlatform
			infrastructure      *fakeinfrastructure.FakeInfrastructure
			blobstore           *fakeblobstore.FakeBlobstore
			taskService         *faketask.FakeService
			notifier            *fakenotif.FakeNotifier
			applier             *fakeappl.FakeApplier
			compiler            *fakecomp.FakeCompiler
			jobSupervisor       *fakejobsuper.FakeJobSupervisor
			specService         *fakeas.FakeV1Service
			drainScriptProvider boshdrain.DrainScriptProvider
			factory             Factory
			logger              boshlog.Logger
		)

		BeforeEach(func() {
			settings = &fakesettings.FakeSettingsService{}
			platform = fakeplatform.NewFakePlatform()
			infrastructure = fakeinfrastructure.NewFakeInfrastructure()
			blobstore = &fakeblobstore.FakeBlobstore{}
			taskService = &faketask.FakeService{}
			notifier = fakenotif.NewFakeNotifier()
			applier = fakeappl.NewFakeApplier()
			compiler = fakecomp.NewFakeCompiler()
			jobSupervisor = fakejobsuper.NewFakeJobSupervisor()
			specService = fakeas.NewFakeV1Service()
			drainScriptProvider = boshdrain.NewConcreteDrainScriptProvider(nil, nil, platform.GetDirProvider())
			logger = boshlog.NewLogger(boshlog.LevelNone)

			factory = NewFactory(
				settings,
				platform,
				infrastructure,
				blobstore,
				taskService,
				notifier,
				applier,
				compiler,
				jobSupervisor,
				specService,
				drainScriptProvider,
				logger,
			)
		})

		It("returns error if action cannot be created", func() {
			action, err := factory.Create("fake-unknown-action")
			Expect(err).To(HaveOccurred())
			Expect(action).To(BeNil())
		})

		It("apply", func() {
			action, err := factory.Create("apply")
			Expect(err).NotTo(HaveOccurred())
			Expect(action).To(Equal(NewApply(applier, specService)))
		})

		It("drain", func() {
			action, err := factory.Create("drain")
			Expect(err).NotTo(HaveOccurred())
			Expect(action).To(Equal(NewDrain(notifier, specService, drainScriptProvider, jobSupervisor)))
		})

		It("fetch_logs", func() {
			action, err := factory.Create("fetch_logs")
			Expect(err).NotTo(HaveOccurred())
			Expect(action).To(Equal(NewLogs(platform.GetCompressor(), platform.GetCopier(), blobstore, platform.GetDirProvider())))
		})

		It("get_task", func() {
			action, err := factory.Create("get_task")
			Expect(err).NotTo(HaveOccurred())
			Expect(action).To(Equal(NewGetTask(taskService)))
		})

		It("get_state", func() {
			ntpService := boshntp.NewConcreteService(platform.GetFs(), platform.GetDirProvider())
			action, err := factory.Create("get_state")
			Expect(err).NotTo(HaveOccurred())
			Expect(action).To(Equal(NewGetState(settings, specService, jobSupervisor, platform.GetVitalsService(), ntpService)))
		})

		It("list_disk", func() {
			action, err := factory.Create("list_disk")
			Expect(err).NotTo(HaveOccurred())
			Expect(action).To(Equal(NewListDisk(settings, platform, logger)))
		})

		It("migrate_disk", func() {
			action, err := factory.Create("migrate_disk")
			Expect(err).NotTo(HaveOccurred())
			Expect(action).To(Equal(NewMigrateDisk(platform, platform.GetDirProvider())))
		})

		It("mount_disk", func() {
			action, err := factory.Create("mount_disk")
			Expect(err).NotTo(HaveOccurred())
			Expect(action).To(Equal(NewMountDisk(settings, infrastructure, platform, platform.GetDirProvider())))
		})

		It("ping", func() {
			action, err := factory.Create("ping")
			Expect(err).NotTo(HaveOccurred())
			Expect(action).To(Equal(NewPing()))
		})

		It("prepare_network_change", func() {
			action, err := factory.Create("prepare_network_change")
			Expect(err).NotTo(HaveOccurred())
			Expect(action).To(Equal(NewPrepareNetworkChange(platform.GetFs(), settings)))
		})

		It("prepare_configure_networks", func() {
			action, err := factory.Create("prepare_configure_networks")
			Expect(err).NotTo(HaveOccurred())
			Expect(action).To(Equal(NewPrepareConfigureNetworks(platform.GetFs(), settings)))
		})

		It("configure_networks", func() {
			action, err := factory.Create("configure_networks")
			Expect(err).NotTo(HaveOccurred())
			Expect(action).To(Equal(NewConfigureNetworks()))
		})

		It("ssh", func() {
			action, err := factory.Create("ssh")
			Expect(err).NotTo(HaveOccurred())
			Expect(action).To(Equal(NewSsh(settings, platform, platform.GetDirProvider())))
		})

		It("start", func() {
			action, err := factory.Create("start")
			Expect(err).NotTo(HaveOccurred())
			Expect(action).To(Equal(NewStart(jobSupervisor)))
		})

		It("stop", func() {
			action, err := factory.Create("start")
			Expect(err).NotTo(HaveOccurred())
			Expect(action).To(Equal(NewStart(jobSupervisor)))
		})

		It("unmount_disk", func() {
			action, err := factory.Create("unmount_disk")
			Expect(err).NotTo(HaveOccurred())
			Expect(action).To(Equal(NewUnmountDisk(settings, platform)))
		})

		It("compile_package", func() {
			action, err := factory.Create("compile_package")
			Expect(err).NotTo(HaveOccurred())
			Expect(action).To(Equal(NewCompilePackage(compiler)))
		})

		It("run_errand", func() {
			action, err := factory.Create("run_errand")
			Expect(err).NotTo(HaveOccurred())
			Expect(action).To(Equal(NewRunErrand(specService, "/var/vcap/jobs", platform.GetRunner())))
		})

		It("prepare", func() {
			action, err := factory.Create("prepare")
			Expect(err).NotTo(HaveOccurred())
			Expect(action).To(Equal(NewPrepare(applier)))
		})
	})
}
