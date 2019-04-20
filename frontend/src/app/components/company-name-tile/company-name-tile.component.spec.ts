import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CompanyNameTileComponent } from './company-name-tile.component';

describe('CompanyNameTileComponent', () => {
  let component: CompanyNameTileComponent;
  let fixture: ComponentFixture<CompanyNameTileComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CompanyNameTileComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CompanyNameTileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
