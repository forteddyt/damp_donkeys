import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SelectCompanyTileComponent } from './select-company-tile.component';

describe('SelectCompanyTileComponent', () => {
  let component: SelectCompanyTileComponent;
  let fixture: ComponentFixture<SelectCompanyTileComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SelectCompanyTileComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SelectCompanyTileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
