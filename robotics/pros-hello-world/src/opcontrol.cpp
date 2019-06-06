#include "main.h"

/**
 * Runs the operator control code. This function will be started in its own task
 * with the default priority and stack size whenever the robot is enabled via
 * the Field Management System or the VEX Competition Switch in the operator
 * control mode.
 *
 * If no competition control is connected, this function will run immediately
 * following initialize().
 *
 * If the robot is disabled or communications is lost, the
 * operator control task will be stopped. Re-enabling the robot will restart the
 * task, not resume it from where it left off.
 */
void opcontrol() {
	pros::Controller master(pros::E_CONTROLLER_MASTER);
	pros::Motor back_left_mtr(1);
	pros::Motor back_right_mtr(2);
	pros::Motor front_right_mtr(10);
	pros::Motor front_left_mtr(9);
	pros::Motor strafe_mtr(8);
	// myMotor.set_encoder_units(pros::E_MOTOR_ENCODER_ROTATIONS);
	while (true) {
		pros::lcd::print(0, "%d %d %d", (pros::lcd::read_buttons() & LCD_BTN_LEFT) >> 2,
		                 (pros::lcd::read_buttons() & LCD_BTN_CENTER) >> 1,
		                 (pros::lcd::read_buttons() & LCD_BTN_RIGHT) >> 0);
		int left = -master.get_analog(ANALOG_LEFT_Y);
		int right = master.get_analog(ANALOG_RIGHT_Y);

		back_left_mtr = left;
		back_right_mtr = right;
		front_left_mtr = left;
		front_right_mtr = right;

		int strafe = master.get_analog(ANALOG_LEFT_X) + master.get_analog(ANALOG_RIGHT_X);
		if (std::abs(strafe) < 64) {
			strafe = 0;
		}
		strafe_mtr = strafe;

		// pros::lcd::print(2, "Current position: %.1d ", myMotor.get_position());
		// int leftButton = pros::lcd::read_buttons() & LCD_BTN_LEFT;
		// int rightButton = pros::lcd::read_buttons() & LCD_BTN_RIGHT;
		// if (leftButton) {
		// 	myMotor.move_absolute(0, 10);
		// }
		// if (rightButton) {
		// 	myMotor.move_absolute(200, -10);
		// }
		// // //myMotor.move(120);



		pros::delay(20);
	}
}
