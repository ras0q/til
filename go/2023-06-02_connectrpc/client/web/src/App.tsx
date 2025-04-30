import { createPromiseClient } from '@bufbuild/connect'
import { createConnectTransport } from '@bufbuild/connect-web'
import './App.css'
import { ReadyService } from './generated/buf/api/proto/ready_connectweb'

function App() {
  const transport = createConnectTransport({
    baseUrl: 'http://localhost:8080/',
  })
  const client = createPromiseClient(ReadyService, transport)

  return (
    <>
      <form
        onSubmit={async (e) => {
          e.preventDefault()
          const res = await client.ready({})
          alert(res)
        }}
      >
        <button type="submit">Ready</button>
      </form>
      <form
        onSubmit={async (e) => {
          e.preventDefault()
          const res = await client.unready({})
          alert(res)
        }}
      >
        <button type="submit">Unready</button>
      </form>
    </>
  )
}

export default App
