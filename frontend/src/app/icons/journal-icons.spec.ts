import {TestBed} from '@angular/core/testing';
import {MatIconRegistry} from '@angular/material/icon';
import {firstValueFrom} from 'rxjs';
import {registerJournalIcons, registerJournalSvgIcons} from './journal-icons';

describe('journal icons', () => {
  it('uses Material Symbols as the default font', () => {
    TestBed.runInInjectionContext(() => registerJournalIcons());

    expect(TestBed.inject(MatIconRegistry).getDefaultFontSetClass()).toEqual(['material-symbols-outlined']);
  });

  it('registers the custom SVG icon without downloading anything', async () => {
    TestBed.runInInjectionContext(() => registerJournalSvgIcons());
    // Registering again is a no-op
    TestBed.runInInjectionContext(() => registerJournalSvgIcons());

    const svg = await firstValueFrom(TestBed.inject(MatIconRegistry).getNamedSvgIcon('checkbox-intermediate-variant'));
    expect(svg.querySelector('path')).not.toBeNull();
  });
});
