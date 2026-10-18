import {inject} from '@angular/core';
import {MatIconRegistry} from '@angular/material/icon';
import {DomSanitizer} from '@angular/platform-browser';

// "Task has an update" status. Not in Material Symbols, so it's the single icon taken from
// Material Design Icons (checkbox-intermediate-variant) instead of downloading the whole set.
const CHECKBOX_INTERMEDIATE = `<svg viewBox="0 0 24 24"><path d="M19 3H5C3.9 3 3 3.9 3 5V19C3 20.1 3.9 21 5 21H19C20.1 21 21 20.1 21 19V5C21 3.9 20.1 3 19 3M19 19H5V5H19V19M7 17V7H17"/></svg>`;

const registered = new WeakSet<MatIconRegistry>();

// App initializer: Material Symbols as the default icon font
export function registerJournalIcons() {
  inject(MatIconRegistry).setDefaultFontSetClass('material-symbols-outlined');
  registerJournalSvgIcons();
}

// Called by the components that use the SVG icons, so they work wherever they're rendered
// (including tests, which don't run app initializers). Registers once per registry.
export function registerJournalSvgIcons() {
  const registry = inject(MatIconRegistry);
  if (registered.has(registry)) {
    return;
  }
  registry.addSvgIconLiteral('checkbox-intermediate-variant', inject(DomSanitizer).bypassSecurityTrustHtml(CHECKBOX_INTERMEDIATE));
  registered.add(registry);
}
