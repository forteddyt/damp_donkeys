import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CompanyWelcomeComponent } from './company-welcome.component';

describe('CompanyWelcomeComponent', () => {
  let component: CompanyWelcomeComponent;
  let fixture: ComponentFixture<CompanyWelcomeComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CompanyWelcomeComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CompanyWelcomeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
