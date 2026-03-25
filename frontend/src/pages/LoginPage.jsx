import React, { useState } from 'react';
import { authApi } from '../api/auth';

export default function LoginPage({ onLogin }) {
    const [id, setId] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    const [showPass, setShowPass] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setLoading(true);
        try {
            const res = await authApi.login(id, password);
            const token = res.data?.data?.token;
            const username = res.data?.data?.username || id;
            if (token) {
                localStorage.setItem('token', token);
                localStorage.setItem('username', username);
                onLogin({ id, username });
            } else {
                setError('Login gagal: token tidak ditemukan.');
            }
        } catch (err) {
            const msg = err.response?.data?.errors || err.response?.statusText || 'Username atau password salah.';
            setError(msg);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div style={styles.wrapper}>
            {/* Animated background blobs */}
            <div style={styles.blob1} />
            <div style={styles.blob2} />

            <div style={styles.card}>
                {/* Logo / branding */}
                <div style={styles.brand}>
                    <div style={styles.iconBox}>
                        <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                            <path d="M22 12h-4l-3 9L9 3l-3 9H2" />
                        </svg>
                    </div>
                    <h1 style={styles.brandTitle}>
                        WAN <span style={{ color: 'var(--accent-glow)' }}>System</span>
                    </h1>
                </div>

                <p style={styles.subtitle}>Masuk ke dashboard monitoring jaringan</p>

                {error && (
                    <div style={styles.errorBox}>
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                            <circle cx="12" cy="12" r="10" /><line x1="12" y1="8" x2="12" y2="12" /><line x1="12" y1="16" x2="12.01" y2="16" />
                        </svg>
                        {error}
                    </div>
                )}

                <form onSubmit={handleSubmit} autoComplete="off">
                    <div style={styles.fieldGroup}>
                        <label style={styles.label}>Username / ID</label>
                        <div style={styles.inputWrap}>
                            <span style={styles.inputIcon}>
                                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" /><circle cx="12" cy="7" r="4" />
                                </svg>
                            </span>
                            <input
                                id="login-id"
                                type="text"
                                placeholder="Masukkan username"
                                value={id}
                                onChange={(e) => setId(e.target.value)}
                                required
                                style={styles.input}
                            />
                        </div>
                    </div>

                    <div style={styles.fieldGroup}>
                        <label style={styles.label}>Password</label>
                        <div style={styles.inputWrap}>
                            <span style={styles.inputIcon}>
                                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                    <rect x="3" y="11" width="18" height="11" rx="2" ry="2" /><path d="M7 11V7a5 5 0 0 1 10 0v4" />
                                </svg>
                            </span>
                            <input
                                id="login-password"
                                type={showPass ? 'text' : 'password'}
                                placeholder="Masukkan password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                required
                                style={styles.input}
                            />
                            <button
                                type="button"
                                onClick={() => setShowPass(!showPass)}
                                style={styles.eyeBtn}
                                tabIndex={-1}
                            >
                                {showPass ? (
                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                        <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94" />
                                        <path d="M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19" />
                                        <line x1="1" y1="1" x2="23" y2="23" />
                                    </svg>
                                ) : (
                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                        <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" /><circle cx="12" cy="12" r="3" />
                                    </svg>
                                )}
                            </button>
                        </div>
                    </div>

                    <button
                        id="login-submit"
                        type="submit"
                        disabled={loading}
                        style={{ ...styles.submitBtn, opacity: loading ? 0.7 : 1 }}
                    >
                        {loading ? (
                            <>
                                <span style={styles.spinner} />
                                Memproses...
                            </>
                        ) : (
                            <>
                                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                    <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4" /><polyline points="10 17 15 12 10 7" /><line x1="15" y1="12" x2="3" y2="12" />
                                </svg>
                                Masuk
                            </>
                        )}
                    </button>
                </form>

                <p style={styles.footer}>WAN Bandwidth Monitoring System &copy; 2024</p>
            </div>

            <style>{`
                @keyframes blob {
                    0%, 100% { transform: translate(0,0) scale(1); }
                    33% { transform: translate(30px,-40px) scale(1.15); }
                    66% { transform: translate(-20px,20px) scale(0.9); }
                }
                @keyframes spin {
                    to { transform: rotate(360deg); }
                }
                #login-id:focus, #login-password:focus {
                    outline: none;
                    border-color: var(--accent-glow) !important;
                    background: rgba(255,255,255,0.1) !important;
                    box-shadow: 0 0 0 3px rgba(34,211,238,0.15);
                }
                #login-submit:hover:not(:disabled) {
                    transform: translateY(-2px);
                    box-shadow: 0 8px 24px rgba(59,130,246,0.6) !important;
                }
            `}</style>
        </div>
    );
}

const styles = {
    wrapper: {
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'var(--bg-color)',
        position: 'relative',
        overflow: 'hidden',
        padding: '20px',
    },
    blob1: {
        position: 'absolute', width: 400, height: 400, borderRadius: '50%',
        background: 'radial-gradient(circle, rgba(59,130,246,0.25), transparent 70%)',
        top: '-100px', left: '-100px',
        animation: 'blob 9s infinite ease-in-out',
    },
    blob2: {
        position: 'absolute', width: 350, height: 350, borderRadius: '50%',
        background: 'radial-gradient(circle, rgba(34,211,238,0.2), transparent 70%)',
        bottom: '-80px', right: '-80px',
        animation: 'blob 12s infinite ease-in-out reverse',
    },
    card: {
        position: 'relative',
        width: '100%',
        maxWidth: 420,
        background: 'rgba(30,41,59,0.75)',
        backdropFilter: 'blur(20px)',
        WebkitBackdropFilter: 'blur(20px)',
        border: '1px solid rgba(255,255,255,0.1)',
        borderRadius: 24,
        padding: '40px 36px',
        boxShadow: '0 25px 60px rgba(0,0,0,0.5)',
    },
    brand: {
        display: 'flex', alignItems: 'center', gap: 12, marginBottom: 8,
    },
    iconBox: {
        width: 44, height: 44, borderRadius: 12,
        background: 'linear-gradient(135deg, #3b82f6, #22d3ee)',
        display: 'flex', alignItems: 'center', justifyContent: 'center',
        color: 'white', boxShadow: '0 4px 14px rgba(59,130,246,0.5)',
    },
    brandTitle: {
        fontSize: '1.6rem', fontWeight: 700, color: 'var(--text-primary)',
    },
    subtitle: {
        color: 'var(--text-secondary)', fontSize: '0.9rem',
        marginBottom: 28, lineHeight: 1.5,
    },
    errorBox: {
        display: 'flex', alignItems: 'center', gap: 8,
        background: 'rgba(239,68,68,0.15)',
        border: '1px solid rgba(239,68,68,0.4)',
        borderRadius: 10, padding: '10px 14px',
        color: '#fca5a5', fontSize: '0.85rem', marginBottom: 20,
    },
    fieldGroup: { marginBottom: 18 },
    label: {
        display: 'block', fontSize: '0.8rem', fontWeight: 600,
        color: 'var(--text-secondary)', marginBottom: 8, letterSpacing: '0.05em',
        textTransform: 'uppercase',
    },
    inputWrap: { position: 'relative', display: 'flex', alignItems: 'center' },
    inputIcon: {
        position: 'absolute', left: 14, color: 'var(--text-secondary)',
        display: 'flex', pointerEvents: 'none',
    },
    input: {
        width: '100%', padding: '12px 44px',
        background: 'rgba(255,255,255,0.05)',
        border: '1px solid rgba(255,255,255,0.1)',
        borderRadius: 10, color: 'white', fontSize: '0.95rem',
        transition: 'all 0.3s ease',
    },
    eyeBtn: {
        position: 'absolute', right: 14,
        background: 'transparent', border: 'none', cursor: 'pointer',
        color: 'var(--text-secondary)', display: 'flex', padding: 4,
    },
    submitBtn: {
        width: '100%', padding: '13px 20px', marginTop: 8,
        background: 'linear-gradient(135deg, #3b82f6, #22d3ee)',
        border: 'none', borderRadius: 10, color: 'white',
        fontSize: '1rem', fontWeight: 700, cursor: 'pointer',
        display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 8,
        boxShadow: '0 4px 14px rgba(59,130,246,0.4)',
        transition: 'all 0.3s ease',
    },
    spinner: {
        width: 18, height: 18, borderRadius: '50%',
        border: '2px solid rgba(255,255,255,0.3)',
        borderTopColor: 'white',
        display: 'inline-block',
        animation: 'spin 0.7s linear infinite',
    },
    footer: {
        marginTop: 28, fontSize: '0.75rem',
        color: 'var(--text-secondary)', textAlign: 'center',
    },
};
