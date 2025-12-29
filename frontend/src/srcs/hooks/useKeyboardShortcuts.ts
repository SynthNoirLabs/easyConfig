import { useEffect } from 'react';

type ShortcutCallback = () => void;
type ShortcutMap = Record<string, ShortcutCallback>;

const getShortcutKey = (e: KeyboardEvent): string => {
  const meta = e.ctrlKey || e.metaKey; // Ctrl on Windows/Linux, Cmd on Mac
  const alt = e.altKey;
  const shift = e.shiftKey;
  const key = e.key.toLowerCase();

  let shortcut = '';
  if (meta) shortcut += 'ctrl+';
  if (alt) shortcut += 'alt+';
  if (shift) shortcut += 'shift+';
  shortcut += key;
  return shortcut;
};

export function useKeyboardShortcuts(shortcuts: ShortcutMap) {
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      const key = getShortcutKey(e);
      if (shortcuts[key]) {
        e.preventDefault();
        shortcuts[key]();
      }
    };

    window.addEventListener('keydown', handler);
    return () => {
      window.removeEventListener('keydown', handler);
    };
  }, [shortcuts]);
}
