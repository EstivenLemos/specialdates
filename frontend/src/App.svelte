<script>
  import { onMount } from 'svelte';

  // Datos
  let dates = [];
  let title = '';
  let description = '';
  let eventAtLocal = ''; // for datetime-local input when adding
  let dateType = '';
  let recurring = false;

  // Modal / edición
  let showModal = false;
  let editing = false;
  let editingItem = null;
  let modalTitle = '';

  // Mensaje de estado
  let statusMsg = '';
  let statusType = ''; // 'success' | 'error' | ''

  // Carga inicial
  async function load() {
    try {
      const res = await fetch('/api/dates');
      dates = await res.json();
    } catch (e) {
      console.error(e);
      dates = [];
      setStatus('Cannot load dates from server', 'error');
    }
  }

  onMount(load);

  // util: ISO <-> datetime-local conversions
  function isoToLocal(iso) {
    if (!iso) return '';
    // remove timezone for input local (browser interprets as local)
    const d = new Date(iso);
    d.setMinutes(d.getMinutes() - d.getTimezoneOffset());
    return d.toISOString().slice(0, 16);
  }

  function localToISO(local) {
    if (!local) return null;
    const d = new Date(local);
    // convert to UTC string iso
    return new Date(d.getTime() - d.getTimezoneOffset() * 60000).toISOString();
  }

  // status helper
  function setStatus(msg, type = 'success', timeout = 3000) {
    statusMsg = msg;
    statusType = type;
    if (timeout) {
      setTimeout(() => {
        statusMsg = '';
        statusType = '';
      }, timeout);
    }
  }

  // Crear
  async function add() {
    if (!title.trim() || !eventAtLocal) {
      setStatus('Title and date are required', 'error');
      return;
    }
    const payload = {
      title: title.trim(),
      description: description.trim(),
      event_at: localToISO(eventAtLocal),
      date_type: dateType.trim(),
      recurring
    };
    try {
      const res = await fetch('/api/dates', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(payload)
      });
      if (!res.ok) throw new Error('Server rejected create');
      const created = await res.json();
      dates = [created, ...dates];
      title = ''; description = ''; eventAtLocal = ''; dateType = ''; recurring = false;
      setStatus('Date created', 'success');
    } catch (err) {
      console.error(err);
      setStatus('Error creating date', 'error');
    }
  }

  // Borrar
  async function del(id) {
    if (!confirm('Delete this date?')) return;
    try {
      const res = await fetch(`/api/dates/${id}`, { method: 'DELETE' });
      if (!res.ok) throw new Error('delete failed');
      dates = dates.filter(d => d.id !== id);
      setStatus('Date deleted', 'success');
    } catch (err) {
      console.error(err);
      setStatus('Error deleting date', 'error');
    }
  }

  // Abrir modal edición
  function openEdit(d) {
    editingItem = { ...d };
    editing = true;
    showModal = true;
    modalTitle = 'Edit Date';
    // set local datetime value in input-friendly format
    editingItem.event_local = isoToLocal(d.event_at);
    // focus management handled in markup
  }

  // Guardar edición
  async function saveEdit() {
    if (!editingItem.title.trim() || !editingItem.event_local) {
      setStatus('Title and date are required', 'error');
      return;
    }
    const payload = {
      title: editingItem.title.trim(),
      description: editingItem.description,
      event_at: localToISO(editingItem.event_local),
      date_type: editingItem.date_type,
      recurring: editingItem.recurring
    };
    try {
      const res = await fetch(`/api/dates/${editingItem.id}`, {
        method: 'PUT',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(payload)
      });
      if (!res.ok) throw new Error('update failed');
      // refresh single item
      const updated = await (await fetch(`/api/dates/${editingItem.id}`)).json();
      dates = dates.map(it => it.id === updated.id ? updated : it);
      showModal = false;
      editing = false;
      editingItem = null;
      setStatus('Date updated', 'success');
    } catch (err) {
      console.error(err);
      setStatus('Error updating date', 'error');
    }
  }

  // Cerrar modal
  function closeModal() {
    showModal = false;
    editing = false;
    editingItem = null;
  }

  // keyboard ESC to close modal
  function onKeydown(e) {
    if (e.key === 'Escape' && showModal) closeModal();
  }

  // small helper for display date
  function displayDate(iso) {
    try {
      return new Date(iso).toLocaleString();
    } catch {
      return iso;
    }
  }
</script>

<svelte:window on:keydown={onKeydown} />

<header class="topbar">
  <div class="container">
    <h1>SpecialDates</h1>
    <p class="tagline">Manage memorable dates – clean, fast and reliable</p>
  </div>
</header>

<main class="container">
  <!-- Left: Add form -->
  <section class="card form-card" aria-labelledby="add-heading">
    <h2 id="add-heading">Add special date</h2>
    <div class="form-grid">
      <label>
        <span>Title</span>
        <input class="input" placeholder="e.g. John's Birthday" bind:value={title} />
      </label>
      <label>
        <span>Description</span>
        <input class="input" placeholder="Optional description" bind:value={description} />
      </label>
      <label>
        <span>Date & time</span>
        <input class="input" type="datetime-local" bind:value={eventAtLocal} />
      </label>
      <label>
        <span>Type</span>
        <input class="input" placeholder="Holiday / Birthday / Anniversary" bind:value={dateType} />
      </label>
      <label class="checkbox-label">
        <input type="checkbox" bind:checked={recurring} />
        <span>Recurring annually</span>
      </label>
    </div>

    <div class="actions">
      <button class="btn primary" on:click={add}>Add Date</button>
      <button class="btn ghost" on:click={() => { title=''; description=''; eventAtLocal=''; dateType=''; recurring=false; }}>Reset</button>
    </div>
  </section>

  <!-- Right: List -->
  <section class="list-section">
    <h2>Upcoming dates</h2>

    {#if dates.length === 0}
      <div class="empty">No special dates found.</div>
    {:else}
      <div class="grid">
        {#each dates as d (d.id)}
          <article class="date-card">
            <div class="card-head">
              <h3 class="card-title">{d.title}</h3>
              <div class="chip">{d.date_type || 'General'}</div>
            </div>
            <p class="card-desc">{d.description}</p>
            <div class="card-meta">
              <div class="datetime">{displayDate(d.event_at)}</div>
              <div class="rec">{d.recurring ? 'Recurring' : 'One-time'}</div>
            </div>
            <div class="card-actions">
              <button class="btn small" on:click={() => openEdit(d)}>Edit</button>
              <button class="btn small danger" on:click={() => del(d.id)}>Delete</button>
            </div>
          </article>
        {/each}
      </div>
    {/if}
  </section>
</main>

<!-- Modal -->
{#if showModal}
  <div class="modal-backdrop" role="dialog" aria-modal="true" aria-label={modalTitle}>
    <div class="modal" tabindex="-1">
      <header class="modal-header">
        <h3>{modalTitle}</h3>
        <button class="icon-btn" aria-label="Close" on:click={closeModal}>✕</button>
      </header>
      <div class="modal-body">
        <label><span>Title</span>
          <input class="input" bind:value={editingItem.title} />
        </label>
        <label><span>Description</span>
          <input class="input" bind:value={editingItem.description} />
        </label>
        <label><span>Date & time</span>
          <input class="input" type="datetime-local" bind:value={editingItem.event_local} />
        </label>
        <label><span>Type</span>
          <input class="input" bind:value={editingItem.date_type} />
        </label>
        <label class="checkbox-label">
          <input type="checkbox" bind:checked={editingItem.recurring} />
          <span>Recurring</span>
        </label>
      </div>
      <footer class="modal-footer">
        <button class="btn primary" on:click={saveEdit}>Save</button>
        <button class="btn ghost" on:click={closeModal}>Cancel</button>
      </footer>
    </div>
  </div>
{/if}

<!-- Status toast -->
{#if statusMsg}
  <div class="toast {statusType}">
    {statusMsg}
  </div>
{/if}

<style>
  :root{
    --bg: #f6f8fb;
    --card: #ffffff;
    --muted: #6b7280;
    --accent: #5b8def;
    --accent-600: #3377ff;
    --danger: #ef4444;
    --glass: rgba(255,255,255,0.6);
    --shadow: 0 8px 20px rgba(16,24,40,0.08);
    --radius: 12px;
    --max-w: 1200px;
  }

  /* layout */
  body {
    margin: 0;
    font-family: Inter, ui-sans-serif, system-ui, -apple-system, "Segoe UI", Roboto, "Helvetica Neue", Arial;
    background: linear-gradient(180deg, var(--bg), #eef3fb 60%);
    color: #0f172a;
    -webkit-font-smoothing:antialiased;
    -moz-osx-font-smoothing:grayscale;
  }

  .container {
    width: 94%;
    max-width: var(--max-w);
    margin: 0 auto;
    padding: 2rem 0;
  }

  .topbar {
    background: linear-gradient(90deg, #ffffff, #f8fbff);
    border-bottom: 1px solid rgba(16,24,40,0.04);
    box-shadow: 0 4px 14px rgba(16,24,40,0.04);
  }
  .topbar .container {
    display:flex;
    align-items:center;
    justify-content:space-between;
    gap:1rem;
    padding:1.25rem 1rem;
  }
  .topbar h1 { margin:0; font-size:1.25rem; letter-spacing:0.2px; }
  .tagline { margin:0; color:var(--muted); font-size:0.95rem; }

  main.container {
    display: grid;
    grid-template-columns: 360px 1fr;
    gap: 1.5rem;
    align-items:start;
  }

  /* cards */
  .card {
    background: var(--card);
    border-radius: var(--radius);
    padding: 1rem 1rem;
    box-shadow: var(--shadow);
    border: 1px solid rgba(15,23,42,0.04);
  }

  .form-card h2 { margin:0 0 0.75rem 0; font-size:1.05rem; color:#0b1220; }
  .form-grid { display:flex; flex-direction:column; gap:0.6rem; }
  label { font-size:0.85rem; color:var(--muted); display:flex; flex-direction:column; gap:0.35rem; }
  .input {
    padding:0.55rem 0.65rem;
    border-radius:8px;
    border:1px solid rgba(15,23,42,0.07);
    background: linear-gradient(180deg, #fff, #fbfdff);
    outline: none;
    font-size:0.95rem;
    box-shadow: inset 0 -1px 0 rgba(16,24,40,0.02);
    transition: box-shadow 120ms ease, border-color 120ms ease, transform 120ms ease;
  }
  .input:focus {
    border-color: var(--accent-600);
    box-shadow: 0 6px 24px rgba(51,119,255,0.10);
    transform: translateY(-1px);
  }

  .checkbox-label { display:flex; align-items:center; gap:0.5rem; font-size:0.95rem; margin-top:0.4rem; color:var(--muted); }

  .actions { display:flex; gap:0.6rem; margin-top:0.8rem; }
  .btn {
    padding: 0.55rem 0.75rem;
    border-radius: 10px;
    border: none;
    cursor: pointer;
    font-weight:600;
    letter-spacing:0.2px;
    transition: transform 120ms ease, box-shadow 120ms ease, background 120ms;
    background: #fff;
    box-shadow: 0 6px 14px rgba(16,24,40,0.04);
  }
  .btn:active { transform: translateY(1px); }

  .btn.primary {
    background: linear-gradient(180deg, var(--accent), var(--accent-600));
    color: white;
    box-shadow: 0 10px 28px rgba(51,119,255,0.18);
  }
  .btn.ghost {
    background: transparent;
    border: 1px solid rgba(15,23,42,0.06);
  }
  .btn.small { padding: 0.35rem 0.6rem; font-size:0.85rem; border-radius:8px; }
  .btn.danger { background: linear-gradient(180deg, #ff6b6b, var(--danger)); color:white; }

  /* list */
  .list-section h2 { margin-top:0; }
  .empty { padding: 1rem; color:var(--muted); background:linear-gradient(180deg,#fbfdff,#fff); border-radius:8px; border:1px dashed rgba(15,23,42,0.03); }

  .grid {
    display:grid;
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
    gap: 1rem;
  }

  .date-card {
    background: linear-gradient(180deg, #fff, #fbfdff);
    padding:0.9rem;
    border-radius: 12px;
    box-shadow: 0 10px 24px rgba(16,24,40,0.06);
    display:flex; flex-direction:column; gap:0.6rem;
    border:1px solid rgba(15,23,42,0.03);
  }
  .card-head { display:flex; justify-content:space-between; align-items:center; gap:0.5rem; }
  .card-title { margin:0; font-size:1.05rem; }
  .chip { font-size:0.78rem; color:var(--muted); background:#f0f6ff; padding:0.3rem 0.45rem; border-radius:999px; }

  .card-desc { color: #334155; font-size:0.95rem; min-height:36px; }
  .card-meta { display:flex; justify-content:space-between; color:var(--muted); font-size:0.85rem; }
  .card-actions { display:flex; gap:0.5rem; margin-top:0.4rem; }

  /* modal */
  .modal-backdrop {
    position: fixed; inset:0; display:flex; align-items:center; justify-content:center;
    background: linear-gradient(rgba(6,8,20,0.35), rgba(6,8,20,0.35));
    padding:1.2rem;
    z-index: 60;
  }
  .modal {
    width: 100%;
    max-width: 720px;
    background: var(--card);
    border-radius: 14px;
    box-shadow: 0 30px 80px rgba(8,15,50,0.3);
    padding: 1rem;
    border:1px solid rgba(15,23,42,0.05);
  }
  .modal-header { display:flex; justify-content:space-between; align-items:center; gap:1rem; }
  .modal-body { display:grid; grid-template-columns: 1fr 1fr; gap:0.7rem; margin-top:0.6rem; }
  .modal-footer { display:flex; justify-content:flex-end; gap:0.6rem; margin-top:1rem; }

  .icon-btn { background:transparent; border:none; font-size:1.05rem; cursor:pointer; }

  /* toast */
  .toast {
    position: fixed; right: 1rem; bottom: 1rem; padding:0.6rem 0.9rem; border-radius:10px; color:white; font-weight:600; z-index: 80;
  }
  .toast.success { background: linear-gradient(90deg,#16a34a,#4ade80); box-shadow:0 8px 28px rgba(16,24,40,0.12); }
  .toast.error { background: linear-gradient(90deg,#ef4444,#fb7185); box-shadow:0 8px 28px rgba(16,24,40,0.12); }

  /* responsive */
  @media (max-width: 900px) {
    main.container { grid-template-columns: 1fr; padding-bottom: 3rem; }
    .modal-body { grid-template-columns: 1fr; }
  }
</style>
