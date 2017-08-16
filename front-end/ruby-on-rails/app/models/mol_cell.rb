class MolCell < ApplicationRecord
	belongs_to :cell
	belongs_to :dataset
end