import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PIDsComponent } from './pids.component';

describe('PIDsComponent', () => {
  let component: PIDsComponent;
  let fixture: ComponentFixture<PIDsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PIDsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PIDsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
