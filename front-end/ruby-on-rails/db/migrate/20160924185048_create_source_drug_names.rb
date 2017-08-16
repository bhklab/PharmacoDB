class CreateSourceDrugNames < ActiveRecord::Migration[5.0]
	def change
		create_table :source_drug_names do |t|
			t.belongs_to :drug, index: true
			t.belongs_to :source, index: true
			t.text :drug_name

			t.timestamps
		end
	end
end
