package common_test

import (
	"os"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/werf/werf/pkg/testing/utils"
)

var werfRepositoryDir string

func init() {
	var err error
	werfRepositoryDir, err = filepath.Abs("../../../")
	if err != nil {
		panic(err)
	}
}

var _ = Describe("context", func() {
	BeforeEach(func() {
		utils.RunSucceedCommand(
			testDirPath,
			"git",
			"clone", werfRepositoryDir, testDirPath,
		)

		utils.RunSucceedCommand(
			testDirPath,
			"git",
			"checkout", "-b", "integration-context-test", "v1.0.10",
		)
	})

	AfterEach(func() {
		utils.RunSucceedCommand(
			testDirPath,
			werfBinPath,
			"purge", "--force",
		)
	})

	type entry struct {
		prepareFixturesFunc   func()
		expectedDigest        string
		expectedDarwinDigest  string
		expectedWindowsDigest string
	}

	var itBody = func(entry entry) {
		entry.prepareFixturesFunc()

		output, err := utils.RunCommand(
			testDirPath,
			werfBinPath,
			"build", "--debug",
		)
		Ω(err).ShouldNot(HaveOccurred())

		if runtime.GOOS == "windows" && entry.expectedWindowsDigest != "" {
			Ω(string(output)).Should(ContainSubstring(entry.expectedWindowsDigest))
		} else if runtime.GOOS == "darwin" && entry.expectedDarwinDigest != "" {
			Ω(string(output)).Should(ContainSubstring(entry.expectedDarwinDigest))
		} else {
			Ω(string(output)).Should(ContainSubstring(entry.expectedDigest))
		}
	}

	var _ = DescribeTable("checksum", itBody,
		Entry("without git", entry{
			prepareFixturesFunc: func() {
				utils.CopyIn(utils.FixturePath("context", "default"), testDirPath)
				Ω(os.RemoveAll(filepath.Join(testDirPath, ".git"))).Should(Succeed())
			},
			expectedDigest:        "ad29bc06ec3893bdfb9beda4720926d2707adea9a9a1b1bdcee9bacd",
			expectedDarwinDigest:  "6419296f73e469ab97cb99defc7dc20c9ad7e9fbf211539e2d0f6639",
			expectedWindowsDigest: "249d3c8d8d2886a030b79fc65d62b921c42e04e5908a40977855e1c5",
		}),
		Entry("with ls-tree", entry{
			prepareFixturesFunc: func() {
				utils.CopyIn(utils.FixturePath("context", "default"), testDirPath)
			},
			expectedDigest:        "70d001f449a48b26160d8b94ab43c48c0209b93a3315236618626f22",
			expectedWindowsDigest: "d5a4acd8b3d55630b0e4d5d9c4cb68467c9ed82a62ed2009243f3119",
		}),
		Entry("with ls-tree and status", entry{
			prepareFixturesFunc: func() {
				utils.RunSucceedCommand(
					testDirPath,
					"git",
					"reset", "HEAD~50",
				)

				utils.CopyIn(utils.FixturePath("context", "default"), testDirPath)
			},
			expectedDigest:        "dc8d845aac3e9894f226c2b816a6c52b477d478d5cf466c5470c86a9",
			expectedWindowsDigest: "b0d542afba27f6b684f8fa2cc1e4f83ccef5ef87e200aa29682dbdf7",
		}),
		Entry("with ls-tree, status and ignored files by .gitignore files", entry{
			prepareFixturesFunc: func() {
				utils.RunSucceedCommand(
					testDirPath,
					"git",
					"reset", "HEAD~50",
				)

				utils.CopyIn(utils.FixturePath("context", "default"), testDirPath)
				utils.CopyIn(utils.FixturePath("context", "gitignores"), testDirPath)
			},
			expectedDigest:        "0ce488f5c941bb516b9c3738fb48f2215069ae64e6d039483e6744ed",
			expectedWindowsDigest: "a89a7669f876e77a98b4f285e2355e5712ab69343332651d68af75bf",
		}),
	)
})
