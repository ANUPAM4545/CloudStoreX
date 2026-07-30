import { format, parseISO, isValid } from "date-fns";
import { ObjectDTO, ObjectViewModel } from "./types";

export function formatBytes(bytes: number, decimals = 1): string {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
}

export function formatTimestamp(ts?: string): string {
  if (!ts) return "-";
  const parsed = parseISO(ts);
  if (!isValid(parsed)) return ts;
  return format(parsed, "MMM d, yyyy • HH:mm");
}

export function buildFolderViewModels(
  objects: ObjectDTO[],
  currentPrefix = ""
): ObjectViewModel[] {
  const foldersMap = new Map<string, ObjectViewModel>();
  const files: ObjectViewModel[] = [];

  for (const obj of objects) {
    if (!obj.key.startsWith(currentPrefix)) continue;

    const relKey = obj.key.slice(currentPrefix.length);
    if (!relKey) continue; // exact prefix match

    const slashIdx = relKey.indexOf("/");
    if (slashIdx !== -1) {
      // It belongs inside a subfolder
      const folderName = relKey.slice(0, slashIdx);
      const folderKey = currentPrefix + folderName + "/";
      if (!foldersMap.has(folderKey)) {
        foldersMap.set(folderKey, {
          key: folderKey,
          name: folderName,
          isFolder: true,
          prefix: folderKey,
          sizeFormatted: "-",
          sizeRaw: 0,
          lastModifiedFormatted: "-",
          contentType: "folder",
        });
      }
    } else {
      // It is a direct file in currentPrefix
      files.push({
        key: obj.key,
        name: relKey,
        isFolder: false,
        prefix: currentPrefix,
        sizeFormatted: formatBytes(obj.size || 0),
        sizeRaw: obj.size || 0,
        lastModifiedFormatted: formatTimestamp(obj.last_modified),
        contentType: obj.content_type || "application/octet-stream",
      });
    }
  }

  const sortedFolders = Array.from(foldersMap.values()).sort((a, b) =>
    a.name.localeCompare(b.name)
  );
  const sortedFiles = files.sort((a, b) => a.name.localeCompare(b.name));

  return [...sortedFolders, ...sortedFiles];
}
