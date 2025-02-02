package cmd

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/miniscruff/changie/core"
)

var _ = Describe("Latest", func() {
	var (
		fs         afero.Fs
		afs        afero.Afero
		testConfig core.Config
	)

	BeforeEach(func() {
		fs = afero.NewMemMapFs()
		afs = afero.Afero{Fs: fs}
		testConfig = core.Config{
			ChangesDir:    "chgs",
			UnreleasedDir: "unrel",
			HeaderPath:    "head.tpl.md",
			ChangelogPath: "changelog.md",
			VersionExt:    "md",
			VersionFormat: "",
			KindFormat:    "",
			ChangeFormat:  "",
			Kinds:         []core.KindConfig{},
		}
		err := testConfig.Save(afs.WriteFile)
		Expect(err).To(BeNil())
	})

	It("echos latest version", func() {
		_, err := afs.Create("chgs/v0.0.1.md")
		Expect(err).To(BeNil())
		_, err = afs.Create("chgs/v0.1.0.md")
		Expect(err).To(BeNil())
		_, err = afs.Create("chgs/head.tpl.md")
		Expect(err).To(BeNil())

		removePrefix = false
		res, err := latestPipeline(afs, false)
		Expect(err).To(BeNil())
		Expect(res).To(Equal("v0.1.0\n"))
	})

	It("echos latest version without prefix", func() {
		_, err := afs.Create("chgs/not-a-version.md")
		Expect(err).To(BeNil())
		_, err = afs.Create("chgs/v0.0.1.md")
		Expect(err).To(BeNil())
		_, err = afs.Create("chgs/v0.1.0.md")
		Expect(err).To(BeNil())

		removePrefix = true
		res, err := latestPipeline(afs, false)
		Expect(err).To(BeNil())
		Expect(res).To(Equal("0.1.0\n"))
	})

	It("fails if bad config file", func() {
		err := afs.Remove(core.ConfigPaths[0])
		Expect(err).To(BeNil())

		_, err = latestPipeline(afs, false)
		Expect(err).NotTo(BeNil())
	})

	It("fails if unable to get versions", func() {
		// no files, means bad read for get versions
		_, err := latestPipeline(afs, false)
		Expect(err).NotTo(BeNil())
	})

	makePrereleases := func(afs afero.Afero) {
		_, err := afs.Create("chgs/not-a-version.md")
		Expect(err).To(BeNil())
		_, err = afs.Create("chgs/v0.0.1.md")
		Expect(err).To(BeNil())
		_, err = afs.Create("chgs/v0.1.0.md")
		Expect(err).To(BeNil())
		_, err = afs.Create("chgs/v0.1.1-rc1.md")
		Expect(err).To(BeNil())
		_, err = afs.Create("chgs/v0.1.1-rc2.md")
		Expect(err).To(BeNil())
		_, err = afs.Create("chgs/v0.2.1-rc3.md")
		Expect(err).To(BeNil())
	}

	It("echos latest version, prereleases not skipped", func() {
		makePrereleases(afs)

		removePrefix = true
		res, err := latestPipeline(afs, false)
		Expect(err).To(BeNil())
		Expect(res).To(Equal("0.2.1-rc3\n"))
	})

	It("echos latest version, prereleases skipped", func() {
		makePrereleases(afs)

		removePrefix = true
		res, err := latestPipeline(afs, true)
		Expect(err).To(BeNil())
		Expect(res).To(Equal("0.1.0\n"))
	})
})
