import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CompanyUniqueComponent } from './company-unique.component';

describe('CompanyUniqueComponent', () => {
  let component: CompanyUniqueComponent;
  let fixture: ComponentFixture<CompanyUniqueComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CompanyUniqueComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CompanyUniqueComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
