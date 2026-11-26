import { useState } from 'react'
import { sendFeedback, type Feedback, type FeedbackSource } from './api'

type Props = {
  problemId: string
  source: FeedbackSource
}

export function FeedbackButtons({ problemId, source }: Props) {
  const [sent, setSent] = useState<Feedback | null>(null)

  const handleClick = (feedback: Feedback) => {
    setSent(feedback)
    sendFeedback(problemId, source, feedback).catch(() => setSent(null))
  }

  if (sent) {
    return <span className="feedback-sent">フィードバックありがとう！</span>
  }

  return (
    <span className="feedback-buttons">
      <button type="button" onClick={() => handleClick('solved')}>
        解けた
      </button>
      <button type="button" onClick={() => handleClick('needed_help')}>
        解けない/解答見た
      </button>
    </span>
  )
}
