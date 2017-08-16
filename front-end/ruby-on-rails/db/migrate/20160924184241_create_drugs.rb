class CreateDrugs < ActiveRecord::Migration[5.0]
	def change
		create_table :drugs, id: false do |t|
			t.integer :drug_id
			t.text :drug_name

			t.timestamps
		end
		execute "ALTER TABLE drugs ADD PRIMARY KEY (drug_id);"
	end
end
