class DrugTarget < ApplicationRecord
	belongs_to :drug
	belongs_to :target
end