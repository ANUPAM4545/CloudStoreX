# CloudStoreX Web Dashboard Architecture

The Web Dashboard (`frontend/`) provides an enterprise cloud storage management UI built with **Next.js 15 (App Router)**, **TypeScript**, **Tailwind CSS**, and **shadcn/ui**.

---

## 1. Feature-Oriented Modular Layout

Business logic is structured into self-contained feature modules under `src/features/`, separating domain logic from reusable UI primitives:

```
src/
├── app/                      # Next.js 15 App Router pages
├── components/ui/            # Presentation primitives (DataTable, PageHeader, ConfirmDialog)
├── features/
│   ├── auth/                 # Enterprise login, JWT session store, and Zod schemas
│   ├── buckets/              # Storage container CRUD, grid/table views, and filtering
│   ├── dashboard/            # Overview statistics and recent bucket cards
│   ├── objects/              # Hierarchical folder tree explorer, breadcrumbs, and stream downloads
│   └── uploads/              # Streaming upload queue, cancellation, and progress widget
├── lib/api/                  # Axios API client with automatic JWT token attachment
└── store/                    # Zustand session and upload queue state
```

---

## 2. Key Capabilities

- **Hierarchical Filesystem Simulation**: Transforms flat S3 object key strings (`docs/2026/report.pdf`) into an interactive hierarchical folder tree (`docs/` → `2026/` → `report.pdf`).
- **Streaming Multipart Uploads**: Upload queue widget (`UploadProgress`) tracks real-time progress across page navigations with `AbortController` cancellation and 100 MB file size bounds checking.
- **Enterprise Design System**: Consumes curated HSL dark-mode palettes, standardized spacing (`--radius: 0.75rem`), and responsive layouts.
