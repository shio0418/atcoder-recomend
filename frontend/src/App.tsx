import { useEffect, useState } from 'react'
import { fetchRecommendations, type RecommendResponse, type ReviewParams } from './api'
import { ReviewList } from './ReviewList'
import { ReviewSettings } from './ReviewSettings'
import { NewCandidateList } from './NewCandidateList'
import { SearchSection } from './SearchSection'
import { CFRecommendations } from './CFRecommendations'
import './App.css'

function App() {
  const [data, setData] = useState<RecommendResponse | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  const loadRecommendations = (params: ReviewParams = {}) => {
    setLoading(true)
    setError(null)
    fetchRecommendations(params)
      .then(setData)
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    loadRecommendations()
  }, [])

  if (!data && loading) {
    return <p>読み込み中...</p>
  }

  if (!data && error) {
    return <p className="error">エラー: {error}</p>
  }

  if (!data) {
    return null
  }

  return (
    <main>
      <h1>AtCoder精進コーチ</h1>

      <section>
        <h2>復習候補</h2>
        <ReviewSettings onApply={loadRecommendations} loading={loading} />
        {error && <p className="error">エラー: {error}</p>}
        <ReviewList problems={data.review} />
      </section>

      <section>
        <h2>新規開拓候補</h2>
        <NewCandidateList problems={data.new} />
      </section>

      <SearchSection />

      <CFRecommendations />
    </main>
  )
}

export default App
