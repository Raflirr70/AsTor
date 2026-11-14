function Logo({ size = 40, className = '' }) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 48 48"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      className={className}
      role="img"
      aria-label="Logo AsmenStore"
    >
      <rect width="48" height="48" rx="12" fill="#2563EB" />
      <path
        d="M24 11 36 17.5v13L24 37 12 30.5v-13L24 11Z"
        stroke="#BFDBFE"
        strokeWidth="2"
        strokeLinejoin="round"
      />
      <path
        d="M12 17.5 24 24l12-6.5M24 24v13"
        stroke="#FFFFFF"
        strokeWidth="2"
        strokeLinejoin="round"
      />
    </svg>
  )
}

export default Logo
