const LOGO_SRC = "/brand/cent-logo.png";

export default function Brand({ className = "" }) {
  return (
    <div className={`brand ${className}`.trim()} aria-label="Cent">
      <span className="brand-symbol" aria-hidden="true">
        <img
          className="brand-logo"
          src={LOGO_SRC}
          alt=""
          onError={(event) => { event.currentTarget.hidden = true; }}
        />
        <span className="brand-fallback">C</span>
      </span>
      <span className="brand-name">Cent</span>
    </div>
  );
}
