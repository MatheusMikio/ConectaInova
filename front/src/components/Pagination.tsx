export function Pagination({ page, count, loading, onChange }: {
  page: number
  count: number
  loading: boolean
  onChange: (page: number) => void
}) {
  return <div className="pagination">
    <button type="button" className="secondary" disabled={page === 1 || loading} onClick={() => onChange(page - 1)}>← Anterior</button>
    <span>Página {page}</span>
    <button type="button" className="secondary" disabled={count < 12 || loading} onClick={() => onChange(page + 1)}>Próxima →</button>
  </div>
}
