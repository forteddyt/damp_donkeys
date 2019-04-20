import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CareerFairTileComponent } from './career-fair-tile.component';

describe('CareerFairTileComponent', () => {
  let component: CareerFairTileComponent;
  let fixture: ComponentFixture<CareerFairTileComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CareerFairTileComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CareerFairTileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
