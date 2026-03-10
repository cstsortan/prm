import { NavLink } from 'react-router-dom'

const linkStyle = ({ isActive }: { isActive: boolean }): React.CSSProperties => ({
  display: 'block',
  padding: '8px 16px',
  borderRadius: '6px',
  textDecoration: 'none',
  color: isActive ? '#fff' : '#cbd5e1',
  backgroundColor: isActive ? '#334155' : 'transparent',
  fontWeight: isActive ? 600 : 400,
  transition: 'background-color 0.15s',
})

export function Nav() {
  return (
    <nav
      style={{
        width: 200,
        minHeight: '100vh',
        backgroundColor: '#1e293b',
        padding: '20px 12px',
        display: 'flex',
        flexDirection: 'column',
        gap: 4,
      }}
    >
      <div
        style={{
          fontSize: 20,
          fontWeight: 700,
          color: '#fff',
          padding: '0 16px 16px',
          borderBottom: '1px solid #334155',
          marginBottom: 12,
        }}
      >
        PRM
      </div>
      <NavLink to="/" style={linkStyle} end>
        Dashboard
      </NavLink>
      <NavLink to="/list" style={linkStyle}>
        List
      </NavLink>
      <NavLink to="/tree" style={linkStyle}>
        Tree
      </NavLink>
      <NavLink to="/search" style={linkStyle}>
        Search
      </NavLink>
    </nav>
  )
}
