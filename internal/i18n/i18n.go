package i18n

import "fmt"

var currentLang = "en"

// translations maps English source strings to language-specific translations.
var translations = map[string]map[string]string{
	"✅ Using repository: %s\n": {
		"jp": "✅ リポジトリ: %s を使用します\n",
	},
	"📥 Fetching pull requests...": {
		"jp": "📥 プルリクエストを取得しています...",
	},
	"📊 Pull Request Statistics": {
		"jp": "📊 プルリクエスト統計",
	},
	"🔢 Basic Metrics:": {
		"jp": "🔢 基本指標:",
	},
	"Metric": {
		"jp": "指標",
	},
	"Value": {
		"jp": "値",
	},
	"Total PRs": {
		"jp": "PR総数",
	},
	"Merged PRs": {
		"jp": "マージ済みPR",
	},
	"WIP PRs": {
		"jp": "WIP PR",
	},
	"Releases (main/master merges)": {
		"jp": "リリース回数（main/masterへのマージ）",
	},
	"Reopened PRs": {
		"jp": "再オープンPR",
	},
	"Reopen Rate": {
		"jp": "再オープン率",
	},
	"Revert-like Merges": {
		"jp": "Revert系マージ数",
	},
	"Hotfix Merges": {
		"jp": "Hotfixマージ数",
	},
	"Hotfix→Release Gap (avg/median)": {
		"jp": "Hotfixと直近リリースの間隔（平均/中央値）",
	},
	"Hotfix w/o prior release": {
		"jp": "直近リリースなしのHotfix",
	},
	"Stability Metrics:": {
		"jp": "🛡️ 安定性指標:",
	},
	"Merge Rate": {
		"jp": "マージ率",
	},
	"⏱️ Timing Metrics:": {
		"jp": "⏱️ 時間指標:",
	},
	"Duration": {
		"jp": "時間",
	},
	"Average Lead Time": {
		"jp": "平均リードタイム",
	},
	"Median Lead Time": {
		"jp": "中央値リードタイム",
	},
	"Average Review Time": {
		"jp": "平均レビュー開始まで",
	},
	"Median Review Time": {
		"jp": "レビュー開始まで（中央値）",
	},
	"Review Time": {
		"jp": "レビュー開始まで",
	},
	"Average Merge Wait Time": {
		"jp": "レビュー後マージ待ち平均",
	},
	"Merge Wait Time": {
		"jp": "レビュー後マージ待ち",
	},
	"Median Merge Wait Time": {
		"jp": "レビュー後マージ待ち中央値",
	},
	"Average Approval→Merge Time": {
		"jp": "承認→マージ平均時間",
	},
	"Approval→Merge Time": {
		"jp": "承認→マージ時間",
	},
	"Median Approval→Merge Time": {
		"jp": "承認→マージ中央値",
	},
	"Reopen→Merge Time": {
		"jp": "再オープン→マージ時間",
	},
	"Lead Time": {
		"jp": "リードタイム",
	},
	"Commit→PR Time": {
		"jp": "コミット→PR時間",
	},
	"Avg Commit→PR Time": {
		"jp": "コミット→PR平均時間",
	},
	"💻 Code Change Metrics:": {
		"jp": "💻 コード変更指標:",
	},
	"Average": {
		"jp": "平均",
	},
	"Files Changed": {
		"jp": "変更ファイル数",
	},
	"Lines Added": {
		"jp": "追加行数",
	},
	"Lines Deleted": {
		"jp": "削除行数",
	},
	"Commits per PR": {
		"jp": "PRあたりコミット数",
	},
	"Commit Frequency/Week": {
		"jp": "週あたりコミット頻度",
	},
	"👥 Collaboration Metrics:": {
		"jp": "👥 コラボレーション指標:",
	},
	"Avg Reviewers per PR": {
		"jp": "PRあたりレビュワー数",
	},
	"Self-Merge Rate": {
		"jp": "セルフマージ率",
	},
	"💬 Code Review Analysis:": {
		"jp": "💬 コードレビュー分析:",
	},
	"Median": {
		"jp": "中央値",
	},
	"Max": {
		"jp": "最大",
	},
	"Review Comments per PR": {
		"jp": "PRあたりレビューコメント",
	},
	"📈 Review Coverage:": {
		"jp": "📈 レビューコメント付与率:",
	},
	"Count": {
		"jp": "件数",
	},
	"Percentage": {
		"jp": "割合",
	},
	"PRs with Review Comments": {
		"jp": "レビューコメントありPR",
	},
	"PRs without Review Comments": {
		"jp": "レビューコメントなしPR",
	},
	"🔍 Review Quality:": {
		"jp": "🔍 レビュー品質:",
	},
	"Review Comment Density": {
		"jp": "コメント密度",
	},
	"%.2f comments/100 lines": {
		"jp": "100行あたりコメント %.2f件",
	},
	"📝 No code review comments found in this period (%d PRs analyzed)": {
		"jp": "📝 この期間にコードレビューコメントはありません (%d 件のPRを解析)",
	},
	"💡 This could indicate:": {
		"jp": "💡 可能性:",
	},
	"   • Code quality is consistently high": {
		"jp": "   • コード品質が安定して高い",
	},
	"   • Team does reviews via other channels": {
		"jp": "   • 別チャネルでレビューしている",
	},
	"   • PRs are small and self-explanatory": {
		"jp": "   • PRが小さく自明",
	},
	"🔀 Merge Type Distribution:": {
		"jp": "🔀 マージ方式の分布:",
	},
	"Merge Type": {
		"jp": "マージ方式",
	},
	"🔧 GitHub Actions Analysis": {
		"jp": "🔧 GitHub Actions 解析",
	},
	"📅 Using default date range: %s to %s\n": {
		"jp": "📅 期間をデフォルトに設定: %s 〜 %s\n",
	},
	"✅ Analyzing repository: %s\n": {
		"jp": "✅ 対象リポジトリ: %s\n",
	},
	"📊 Period: %s to %s\n": {
		"jp": "📊 期間: %s 〜 %s\n",
	},
	"🔄 Fetching workflow runs...": {
		"jp": "🔄 ワークフロー実行履歴を取得しています...",
	},
	"⚠️  No workflow runs found in the specified period": {
		"jp": "⚠️  指定期間のワークフロー実行はありません",
	},
	"🎯 GitHub Actions Analytics": {
		"jp": "🎯 GitHub Actions 分析",
	},
	"📊 Summary Statistics:": {
		"jp": "📊 サマリー:",
	},
	"Total Runs": {
		"jp": "実行数",
	},
	"Successful Runs": {
		"jp": "成功",
	},
	"Failed Runs": {
		"jp": "失敗",
	},
	"Success Rate": {
		"jp": "成功率",
	},
	"Avg Duration": {
		"jp": "平均時間",
	},
	"🔄 Workflow Breakdown:": {
		"jp": "🔄 ワークフロー別内訳:",
	},
	"Workflow": {
		"jp": "ワークフロー",
	},
	"Runs": {
		"jp": "実行",
	},
	"Success": {
		"jp": "成功",
	},
	"Failed": {
		"jp": "失敗",
	},
	"⚡ Trigger Event Analysis:": {
		"jp": "⚡ トリガーイベント分析:",
	},
	"Event": {
		"jp": "イベント",
	},
	"❌ Failure Analysis:": {
		"jp": "❌ 失敗解析:",
	},
	"🔴 Failure #%d:": {
		"jp": "🔴 失敗 #%d:",
	},
	"  Workflow: %s\n": {
		"jp": "  ワークフロー: %s\n",
	},
	"  Run: %s\n": {
		"jp": "  実行: %s\n",
	},
	"  Date: %s\n": {
		"jp": "  日時: %s\n",
	},
	"  Duration: %s\n": {
		"jp": "  所要時間: %s\n",
	},
	"  Failed Job: %s\n": {
		"jp": "  失敗ジョブ: %s\n",
	},
	"  Failed Step: %s\n": {
		"jp": "  失敗ステップ: %s\n",
	},
	"  URL: %s\n": {
		"jp": "  URL: %s\n",
	},
	"\n... and %d more failures\n": {
		"jp": "\n...さらに %d 件の失敗があります\n",
	},
}

// SetLanguage configures the output language. Unknown values fall back to English.
func SetLanguage(lang string) {
	if lang == "jp" {
		currentLang = "jp"
		return
	}
	currentLang = "en"
}

// Lang returns the currently configured language.
func Lang() string {
	return currentLang
}

// T returns the translated message if available.
func T(msg string) string {
	if currentLang == "en" {
		return msg
	}
	if m, ok := translations[msg]; ok {
		if t, ok := m[currentLang]; ok && t != "" {
			return t
		}
	}
	return msg
}

// Sprintf formats a translated string with the provided arguments.
func Sprintf(msg string, args ...interface{}) string {
	return fmt.Sprintf(T(msg), args...)
}
