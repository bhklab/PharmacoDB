class CreateExperiments < ActiveRecord::Migration[5.0]
	def change
		create_table :experiments, id: false do |t|
			t.integer :experiment_id
			t.belongs_to :cell, index: true
			t.belongs_to :drug, index: true
			t.belongs_to :dataset, index: true

			t.timestamps
		end
		execute "ALTER TABLE experiments ADD PRIMARY KEY (experiment_id);"
	end
end
