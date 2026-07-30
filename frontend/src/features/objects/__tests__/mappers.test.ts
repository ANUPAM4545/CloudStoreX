import { describe, it, expect } from "vitest";
import { formatBytes, buildFolderViewModels } from "../mappers";

describe("formatBytes", () => {
  it("should format bytes to human-readable units", () => {
    expect(formatBytes(0)).toBe("0 B");
    expect(formatBytes(1024)).toBe("1 KB");
    expect(formatBytes(1536)).toBe("1.5 KB");
    expect(formatBytes(1048576)).toBe("1 MB");
  });
});

describe("buildFolderViewModels", () => {
  it("should split keys into subfolders and direct files", () => {
    const dtos = [
      { key: "docs/2026/report.pdf", bucket: "my-bucket", size: 2048 },
      { key: "docs/notes.txt", bucket: "my-bucket", size: 512 },
      { key: "readme.md", bucket: "my-bucket", size: 100 },
    ];

    const rootItems = buildFolderViewModels(dtos, "");
    expect(rootItems).toHaveLength(2); // "docs/" folder + "readme.md" file
    expect(rootItems[0].name).toBe("docs");
    expect(rootItems[0].isFolder).toBe(true);
    expect(rootItems[1].name).toBe("readme.md");
    expect(rootItems[1].isFolder).toBe(false);

    const docsItems = buildFolderViewModels(dtos, "docs/");
    expect(docsItems).toHaveLength(2); // "2026/" folder + "notes.txt" file
    expect(docsItems[0].name).toBe("2026");
    expect(docsItems[0].isFolder).toBe(true);
    expect(docsItems[1].name).toBe("notes.txt");
    expect(docsItems[1].isFolder).toBe(false);
  });
});
