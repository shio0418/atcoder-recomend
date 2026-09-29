import { useEffect, useState } from 'react'
import { fetchRecommendations, type RecommendResponse } from './api'
import { ReviewList } from './ReviewList'
import { NewCandidateList } from './NewCandidateList'
import './App.css'

function App() {
  const [data, setData] = useState<RecommendResponse | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchRecommendations()
      .then(setData)
      .catch((e: Error) => setError(e.message))
  }, [])

  if (error) {
    return <p className="error">エラー: {error}</p>
  }

  if (!data) {
    return <p>読み込み中...</p>
  }

  return (
    <main>
      <h1>AtCoder精進コーチ</h1>

      <section>
        <h2>復習候補</h2>
        <ReviewList problems={data.review} />
      </section>

      <section>
        <h2>新規開拓候補</h2>
        <NewCandidateList problems={data.new} />
      </section>
    </main>
  )
}

export default App
