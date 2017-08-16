class CreateDoseResponses < ActiveRecord::Migration[5.0]
	def change
		create_table :dose_responses do |t|
			t.belongs_to :experiment, index: true
			t.float :dose
			t.float :response

			t.timestamps
		end
	end
end
